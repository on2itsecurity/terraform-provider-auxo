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
	_ datasource.DataSource              = &protectsurfaceDataSource{}
	_ datasource.DataSourceWithConfigure = &protectsurfaceDataSource{}
)

type protectsurfaceDataSource struct {
	client *auxo.Client
}

type protectsurfaceDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	Uniqueness_key types.String `tfsdk:"uniqueness_key"`
	Name           types.String `tfsdk:"name"`
	//TODO Add all fields
}

// NewprotectsurfaceDataSource is a helper function to simplify the provider implementation.
func NewProtectsurfaceDataSource() datasource.DataSource {
	return &protectsurfaceDataSource{}
}

// Metadata returns the data source type name.
func (d *protectsurfaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_protectsurface"
}

func (d *protectsurfaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	// Retrieve the client from the provider config
	d.client = req.ProviderData.(*auxo.Client)
}

// Schema defines the schema for the data source.
func (d *protectsurfaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "A protectsurface which can be used in a state object to reflect the protectsurface of the resources. Either by specifying the name or the uniqueness_key",
		MarkdownDescription: "A protectsurface which can be used in a state object to reflect the protectsurface of the resources. Either by specifying the name or the uniqueness_key",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Computed unique IDs of the protectsurface",
				MarkdownDescription: "Computed unique IDs of the protectsurface",
				Computed:            true,
			},
			"uniqueness_key": schema.StringAttribute{
				Description:         "Uniqueness key of the protectsurface",
				MarkdownDescription: "Uniqueness key of the protectsurface",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the protectsurface",
				MarkdownDescription: "Name of the protectsurface",
				Optional:            true,
			},
			//TODO Add all fields
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *protectsurfaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state protectsurfaceDataSourceModel

	//Get input
	var input protectsurfaceDataSourceModel
	diags := req.Config.Get(ctx, &input)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Get protectsurfaces
	protectsurfaces, err := d.client.ZeroTrust.GetProtectSurfaces()
	if err != nil {
		resp.Diagnostics.AddError("Unable to retrieve protectsurfaces", err.Error())
		return
	}

	//Make map to check for duplicates (which allowed, but makes it impossible to look for a protectsurface by name)
	psCount := make(map[string]int)
	for _, protectsurface := range protectsurfaces {
		psCount[protectsurface.Name]++
	}
	if psCount[input.Name.ValueString()] > 1 {
		resp.Diagnostics.AddError("Duplicate name on backend, please use uniqueness key", "Duplicate name on backend name "+input.Name.ValueString()+", please use uniqueness key")
		return
	}

	//Check if one of the two is set
	if (input.Uniqueness_key.IsNull() && input.Name.IsNull()) || (!input.Uniqueness_key.IsNull() && !input.Name.IsNull()) {
		resp.Diagnostics.AddError("Either uniqueness_key OR name must be set", "")
		return
	}

	//Find the protectsurface
	for _, l := range protectsurfaces {
		if (l.UniquenessKey == input.Uniqueness_key.ValueString()) || (l.Name == input.Name.ValueString()) {
			state.ID = types.StringValue(l.ID)
			state.Name = types.StringValue(l.Name)
			state.Uniqueness_key = types.StringValue(l.UniquenessKey)
			//TODO Add all fields
			break
		}
	}

	if state.ID.IsNull() {
		resp.Diagnostics.AddError("Unable to find protectsurface", "Unable to find protectsurface with uniqueness_key "+input.Uniqueness_key.ValueString()+" or name "+input.Name.ValueString())
		return
	}

	//set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
