package auxo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceContact() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContactRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceContactRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	provider := m.(*AuxoProvider)
	apiClient := provider.APIClient
	contacts, err := apiClient.CRM.GetContacts()

	if err != nil {
		return diag.FromErr(err)
	}

	for _, c := range contacts {
		if c.Email == d.Get("email").(string) {
			d.SetId(c.ID)
			d.Set("id", c.ID)
			d.Set("email", c.Email)
			break
		}
	}

	if d.Id() == "" {
		return diag.Errorf("Contact not found [%s]", d.Get("email").(string))
	}

	return diags
}
