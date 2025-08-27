package auxo

import (
	"context"
	"os"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/on2itsecurity/go-auxo/v2"
)

// Ensure the implementation satisfies the provider.Provider interface.
var _ provider.Provider = (*auxoProvider)(nil)

// auxoProvider represents the provider.Provider interface
type auxoProvider struct{}

type auxoProviderModel struct {
	Url    types.String `tfsdk:"url"`
	Token  types.String `tfsdk:"token"`
	Name   types.String `tfsdk:"name"`
	Config types.String `tfsdk:"config"`
}

type auxoClient struct {
	client *auxo.Client
	m      *sync.Mutex
}

// New returns a new provider.Provider.
func New() provider.Provider {
	return &auxoProvider{}
}

func (p *auxoProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "auxo"
	resp.Version = "dev"
}

func (p *auxoProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The alias in the ztctl configuration file, this takes precedence over the url and token attributes",
				Description:         "The alias in the ztctl configuration file, this takes precedence over the url and token attributes",
			},
			"config": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Location of the ztctl configuration file, will default to `~/.ztctl/config.json`",
				Description:         "Location of the ztctl configuration file, will default to `~/.ztctl/config.json`",
			},
			"url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The URL of the Auxo API",
				Description:         "The URL of the Auxo API",
			},
			"token": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The token to access the API",
				Description:         "The token to access the API",
				Sensitive:           true,
			},
		},
	}
}

func (p *auxoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	//Default or environment variables
	// checkov:skip=CKV_SECRET_6: False/Positive
	token := os.Getenv("AUXO_TOKEN")
	url := "api.on2it.net"
	config := getDefaultConfigLocation()

	var data auxoProviderModel

	// Read configuration data into model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Check configuration data, which should take precedence over
	// environment variable data, if found.
	if data.Token.ValueString() != "" {
		token = data.Token.ValueString()
	}

	if data.Url.ValueString() != "" {
		url = data.Url.ValueString()
	}

	if data.Config.ValueString() != "" {
		config = data.Config.ValueString()
	}

	if data.Name.ValueString() != "" {
		alias := data.Name.ValueString()

		//Read configuration and set url and token
		cfg, err := getConfig(config, alias)

		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to read configuration file",
				"An unexpected error occurred when reading the configuration file. "+
					"client error: "+err.Error())
		}

		//Find specific alias and set url and token
		url = cfg.APIAddress
		token = cfg.Token
	}

	//Error checking
	if token == "" {
		resp.Diagnostics.AddError(
			"Missing API Token Configuration",
			"While configuring the provider, the API token was not found in "+
				"the AUXO_TOKEN environment variable or provider "+
				"configuration block 'token' or 'config' attribute.",
		)
	}

	if url == "" {
		resp.Diagnostics.AddError(
			"Missing API URL Configuration",
			"While configuring the provider, the API URL was not found in "+
				"the provider configuration block 'url' or 'config' attribute.",
		)
	}

	// Create data/clients and persist to resp.DataSourceData and
	// resp.ResourceData as appropriate.
	client, err := auxo.NewClient(url, token, false)
	c := &auxoClient{
		client: client,
		m:      &sync.Mutex{},
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to create AUXO API client",
			"An unexpected error occurred when creating the AUXO API client. "+
				"client error: "+err.Error())
	}

	resp.DataSourceData = client
	resp.ResourceData = c
}

func (p *auxoProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProtectsurfaceResource,
		NewLocationResource,
		NewMeasureResource,
		NewStateResource,
		NewTransactionflowResource,
	}
}

func (p *auxoProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAssetDataSource,
		NewContactDataSource,
		NewLocationDataSource,
		NewProtectsurfaceDataSource,
	}
}
