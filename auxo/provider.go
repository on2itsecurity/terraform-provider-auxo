package auxo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/on2itsecurity/go-auxo"
)

// AuxoProvider represents the Auxo provider
type AuxoProvider struct {
	APIClient *auxo.Client
}

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "api.on2it.net",
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("AUXOTOKEN", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"auxo_location":       resourceLocation(),
			"auxo_protectsurface": resourceProtectSurface(),
			"auxo_state":          resourceState(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	url := d.Get("url").(string)
	token := d.Get("token").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if token != "" {
		client, err := auxo.NewClient(url, token, false)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		provider := new(AuxoProvider)
		provider.APIClient = client

		return provider, diags
	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "No token specified.",
		Detail:   "Please set the token in the configuration or as environment variable AUXOTOKEN.",
	})

	return nil, diags //Will never be hit
}
