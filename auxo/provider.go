package auxo

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/on2itsecurity/go-auxo"
)

// Ensure the implementation satisfies the provider.Provider interface.
var _ provider.Provider = (*auxoProvider)(nil)

// auxoProvider represents the provider.Provider interface
type auxoProvider struct{}

type auxoProviderModel struct {
	Url    types.String `tfsdk:"url"`
	Token  types.String `tfsdk:"token"`
	Config types.String `tfsdk:"config"`
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
			"config": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The config of the ztctl configuration, this takes precedence over the url and token attributes",
				Description:         "The config of the ztctl configuration, this takes precedence over the url and token attributes",
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
	token := os.Getenv("AUXO_TOKEN")
	url := "api.on2it.net"

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
		//Read configuration and set url and token
		//TODO Read config configuration and set variables
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
	auxoClient, err := auxo.NewClient(url, token, false)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create AUXO API client",
			"An unexpected error occurred when creating the AUXO API client. "+
				"client error: "+err.Error())
	}

	resp.DataSourceData = auxoClient
	resp.ResourceData = auxoClient
}

func (p *auxoProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProtectsurfaceResource,
		NewLocationResource,
		NewStateResource,
	}
}

func (p *auxoProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewcontactDataSource,
	}
}

// // AuxoProvider represents the Auxo provider
// type AuxoProvider struct {
// 	APIClient *auxo.Client
// }

// // Provider -
// func Provider() *schema.Provider {
// 	return &schema.Provider{
// 		Schema: map[string]*schema.Schema{
// 			"url": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Default:  "api.on2it.net",
// 			},
// 			"token": {
// 				Type:        schema.TypeString,
// 				Required:    true,
// 				Sensitive:   true,
// 				DefaultFunc: schema.EnvDefaultFunc("AUXOTOKEN", ""),
// 			},
// 		},
// 		DataSourcesMap: map[string]*schema.Resource{
// 			"auxo_contact": dataSourceContact(),
// 		},
// 		ResourcesMap: map[string]*schema.Resource{
// 			"auxo_location":        resourceLocation(),
// 			"auxo_protectsurface":  resourceProtectSurface(),
// 			"auxo_state":           resourceState(),
// 			"auxo_transactionflow": resourceTransactionFlow(),
// 		},
// 		ConfigureContextFunc: providerConfigure,
// 	}
// }

// func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
// 	url := d.Get("url").(string)
// 	token := d.Get("token").(string)

// 	// Warning or errors can be collected in a slice type
// 	var diags diag.Diagnostics

// 	if token != "" {
// 		client, err := auxo.NewClient(url, token, false)
// 		if err != nil {
// 			return nil, diag.FromErr(err)
// 		}

// 		provider := new(AuxoProvider)
// 		provider.APIClient = client

// 		return provider, diags
// 	}

// 	diags = append(diags, diag.Diagnostic{
// 		Severity: diag.Error,
// 		Summary:  "No token specified.",
// 		Detail:   "Please set the token in the configuration or as environment variable AUXOTOKEN.",
// 	})

// 	return nil, diags //Will never be hit
// }
