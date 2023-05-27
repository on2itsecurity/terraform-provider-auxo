package auxo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/on2itsecurity/go-auxo"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &contactDataSource{}
	_ datasource.DataSourceWithConfigure = &contactDataSource{}
)

type contactDataSource struct {
	client *auxo.Client
}

type contactDataSourceModel struct {
	ID    types.String `tfsdk:"id"`
	Email types.String `tfsdk:"email"`
}

// NewcontactDataSource is a helper function to simplify the provider implementation.
func NewcontactDataSource() datasource.DataSource {
	return &contactDataSource{}
}

// Metadata returns the data source type name.
func (d *contactDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contact"
}

func (d *contactDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	// Retrieve the client from the provider config
	d.client = req.ProviderData.(*auxo.Client)
}

// Schema defines the schema for the data source.
func (d *contactDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "A contact which can be used a.o. as main- or securitycontact in a `protectsurface`.",
		MarkdownDescription: "A contact which can be used a.o. as main- or securitycontact in a `protectsurface`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Computed unique IDs of the contact",
				MarkdownDescription: "Computed unique IDs of the contact",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				Description:         "Emails of the contact",
				MarkdownDescription: "Emails of the contact",
				Required:            true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *contactDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state contactDataSourceModel

	//Get contacts
	contacts, err := d.client.CRM.GetContacts()
	if err != nil {
		resp.Diagnostics.AddError("Unable to retrieve contacts", err.Error())
		return
	}

	//Get input
	var input contactDataSourceModel
	diags := req.Config.Get(ctx, &input)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Find the contact
	for _, c := range contacts {
		if c.Email == input.Email.ValueString() {
			state.ID = types.StringValue(c.ID)
			state.Email = types.StringValue(c.Email)
			break
		}
	}

	//set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
