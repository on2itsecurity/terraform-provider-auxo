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
	_ datasource.DataSource              = &assetDataSource{}
	_ datasource.DataSourceWithConfigure = &assetDataSource{}
)

type assetDataSource struct {
	client *auxo.Client
}

type assetDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// NewassetDataSource is a helper function to simplify the provider implementation.
func NewAssetDataSource() datasource.DataSource {
	return &assetDataSource{}
}

// Metadata returns the data source type name.
func (d *assetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asset"
}

func (d *assetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	// Retrieve the client from the provider config
	d.client = req.ProviderData.(*auxo.Client)
}

// Schema defines the schema for the data source.
func (d *assetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "A asset which can be used as an `exists_on_asset` in a `state` resource.",
		MarkdownDescription: "A asset which can be used as an `exists_on_asset` in a `state` resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Computed unique IDs of the asset",
				MarkdownDescription: "Computed unique IDs of the asset",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the asset",
				MarkdownDescription: "Name of the asset",
				Required:            true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *assetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state assetDataSourceModel

	//Get assets
	assets, err := d.client.Asset.GetAssets()
	if err != nil {
		resp.Diagnostics.AddError("Unable to retrieve assets", err.Error())
		return
	}

	//Get input
	var input assetDataSourceModel
	diags := req.Config.Get(ctx, &input)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Find the asset
	for _, a := range assets {
		if a.Name == input.Name.ValueString() {
			state.ID = types.StringValue(a.ID)
			state.Name = types.StringValue(a.Name)
			break
		}
	}

	if state.ID.IsNull() {
		resp.Diagnostics.AddError("Unable to find asset", "Unable to find asset with name "+input.Name.ValueString())
		return
	}

	//set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
