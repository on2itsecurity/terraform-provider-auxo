package auxo

import (
	"context"

	"github.com/on2itsecurity/go-auxo/zerotrust"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLocation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLocationCreate,
		ReadContext:   resourceLocationRead,
		UpdateContext: resourceLocationUpdate,
		DeleteContext: resourceLocationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the resource/location",
			},
			"uniqueness_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Unique key to generate the ID - only needed for parallel import",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique name of the location",
			},
			"latitude": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "Latitude of the location",
			},
			"longitude": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "Longitude of the location",
			},
		},
	}
}

func resourceLocationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	provider := m.(*AuxoProvider)
	apiClient := provider.APIClient

	var location = new(zerotrust.Location)

	location.UniquenessKey = d.Get("uniqueness_key").(string)
	location.Name = d.Get("name").(string)
	location.Coords.Latitude = d.Get("latitude").(float64)
	location.Coords.Longitude = d.Get("longitude").(float64)

	result, err := apiClient.ZeroTrust.CreateLocationByObject(*location, false)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.ID)

	resourceLocationRead(ctx, d, m)

	return diags
}

func resourceLocationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	provider := m.(*AuxoProvider)
	apiClient := provider.APIClient

	location, err := apiClient.ZeroTrust.GetLocationByID(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("id", location.ID)
	d.Set("uniqueness_key", location.UniquenessKey)
	d.Set("name", location.Name)
	d.Set("latitude", location.Coords.Latitude)
	d.Set("longitude", location.Coords.Longitude)

	return diags
}

func resourceLocationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*AuxoProvider)
	apiClient := provider.APIClient

	location, err := apiClient.ZeroTrust.GetLocationByID(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("uniqueness_key") {
		location.UniquenessKey = d.Get("uniqueness_key").(string)
	}
	if d.HasChange("name") {
		location.Name = d.Get("name").(string)
	}
	if d.HasChange("latitude") {
		location.Coords.Latitude = d.Get("latitude").(float64)
	}
	if d.HasChange("longitude") {
		location.Coords.Longitude = d.Get("longitude").(float64)
	}

	_, err = apiClient.ZeroTrust.UpdateLocation(*location)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceLocationRead(ctx, d, m)
}

func resourceLocationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*AuxoProvider)
	apiClient := provider.APIClient

	var diags diag.Diagnostics

	err := apiClient.ZeroTrust.DeleteLocationByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
