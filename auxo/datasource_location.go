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
	_ datasource.DataSource              = &locationDataSource{}
	_ datasource.DataSourceWithConfigure = &locationDataSource{}
)

type locationDataSource struct {
	client *auxo.Client
}

type locationDataSourceModel struct {
	ID             types.String  `tfsdk:"id"`
	Uniqueness_key types.String  `tfsdk:"uniqueness_key"`
	Name           types.String  `tfsdk:"name"`
	Latitude       types.Float64 `tfsdk:"latitude"`
	Longitude      types.Float64 `tfsdk:"longitude"`
}

// NewLocationDataSource is a helper function to simplify the provider implementation.
func NewLocationDataSource() datasource.DataSource {
	return &locationDataSource{}
}

// Metadata returns the data source type name.
func (d *locationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_location"
}

func (d *locationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	// Retrieve the client from the provider config
	d.client = req.ProviderData.(*auxo.Client)
}

// Schema defines the schema for the data source.
func (d *locationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "A location which can be used in a state object to reflect the location of the resources.",
		MarkdownDescription: "A location which can be used in a state object to reflect the location of the resources",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Computed unique IDs of the location",
				MarkdownDescription: "Computed unique IDs of the location",
				Computed:            true,
			},
			"uniqueness_key": schema.StringAttribute{
				Description:         "Uniqueness key of the location",
				MarkdownDescription: "Uniqueness key of the location",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the location",
				MarkdownDescription: "Name of the location",
				Optional:            true,
			},
			"latitude": schema.Float64Attribute{
				Description:         "Latitude of the resource location",
				MarkdownDescription: "Latitude of the resource location",
				Computed:            true,
			},
			"longitude": schema.Float64Attribute{
				Description:         "Longitude of the resource location",
				MarkdownDescription: "Longitude of the resource location",
				Computed:            true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *locationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state locationDataSourceModel

	//Get input
	var input locationDataSourceModel
	diags := req.Config.Get(ctx, &input)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Get locations
	locations, err := d.client.ZeroTrust.GetLocations()
	if err != nil {
		resp.Diagnostics.AddError("Unable to retrieve locations", err.Error())
		return
	}

	//Check if one of the two is set
	if (input.Uniqueness_key.IsNull() && input.Name.IsNull()) || (!input.Uniqueness_key.IsNull() && !input.Name.IsNull()) {
		resp.Diagnostics.AddError("Either uniqueness_key OR name must be set", "")
		return
	}

	//Find the location //TODO what if multiple entries have the samen name ?
	for _, l := range locations {
		if (l.UniquenessKey == input.Uniqueness_key.ValueString()) || (l.Name == input.Name.ValueString()) {
			state.ID = types.StringValue(l.ID)
			state.Name = types.StringValue(l.Name)
			state.Uniqueness_key = types.StringValue(l.UniquenessKey)
			state.Latitude = types.Float64Value(l.Coords.Latitude)
			state.Longitude = types.Float64Value(l.Coords.Longitude)
			break
		}
	}

	if state.ID.IsNull() {
		resp.Diagnostics.AddError("Unable to find location", "Unable to find location with uniqueness_key "+input.Uniqueness_key.ValueString()+" or name "+input.Name.ValueString())
		return
	}

	//set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
