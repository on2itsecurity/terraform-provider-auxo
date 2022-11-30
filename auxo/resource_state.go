package auxo

import (
	"context"

	"github.com/on2itsecurity/go-auxo/zerotrust"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceState() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStateCreate,
		ReadContext:   resourceStateRead,
		UpdateContext: resourceStateUpdate,
		DeleteContext: resourceStateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the resource/state",
			},
			"uniqueness_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Unique key to generate the ID - only needed for parallel import",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the segment",
			},
			"protectsurface_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ProtectSurface ID",
				ForceNew:    true,
			},
			"location_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Location ID",
				ForceNew:    true,
			},
			"content_type": {
				Type:        schema.TypeString,
				Default:     "static_ipv4",
				Description: "Content type of the state i.e. static_ipv4, static_ipv6, azure_resource",
				Optional:    true,
			},
			"exists_on_assets": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Contains asset IDs which could match this state",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"maintainer": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "api",
				Description: "Maintainer of the state either api or portal_manual",
			},
			"content": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Content of the state e.g. \"10.1.1.2/32\",\"10.1.1.3/32\"",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceStateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	provider := m.(*AuxoProvider)
	apiClient := provider.APIClient

	var state = new(zerotrust.State)

	state.UniquenessKey = d.Get("uniqueness_key").(string)
	state.Description = d.Get("description").(string)
	state.ProtectSurface = d.Get("protectsurface_id").(string)
	state.Location = d.Get("location_id").(string)
	state.ContentType = d.Get("content_type").(string)
	state.Maintainer = d.Get("maintainer").(string)
	state.ExistsOnAssetIDs = createStringSliceFromListInput(d.Get("exists_on_assets").(*schema.Set).List())
	content := createStringSliceFromListInput(d.Get("content").(*schema.Set).List())
	state.Content = &content

	result, err := apiClient.ZeroTrust.CreateStateByObject(*state)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.ID)

	resourceStateRead(ctx, d, m)

	return diags
}

func resourceStateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	provider := m.(*AuxoProvider)
	apiClient := provider.APIClient

	state, err := apiClient.ZeroTrust.GetStateByID(d.Id())

	if err != nil {
		apiError := getAPIError(err)

		//NotExists
		if apiError.ID == "410" {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	d.Set("id", state.ID)
	d.Set("uniqueness_key", state.UniquenessKey)
	d.Set("description", state.Description)
	d.Set("protectsurface_id", state.ProtectSurface)
	d.Set("location_id", state.Location)
	d.Set("content_type", state.ContentType)
	d.Set("maintainer", state.Maintainer)
	d.Set("exists_on_assets", state.ExistsOnAssetIDs)
	d.Set("content", state.Content)

	return diags
}

func resourceStateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	provider := m.(*AuxoProvider)
	apiClient := provider.APIClient

	state, err := apiClient.ZeroTrust.GetStateByID(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("uniqueness_key") {
		state.UniquenessKey = d.Get("uniqueness_key").(string)
	}
	if d.HasChange("description") {
		state.Description = d.Get("description").(string)
	}
	if d.HasChange("protectsurface_id") {
		state.ProtectSurface = d.Get("protectsurface_id").(string)
	}
	if d.HasChange("location_id") {
		state.Location = d.Get("location_id").(string)
	}
	if d.HasChange("content_type") {
		state.ContentType = d.Get("content_type").(string)
	}
	if d.HasChange("maintainer") {
		state.Maintainer = d.Get("maintainer").(string)
	}
	if d.HasChange("exists_on_assets") {
		state.ExistsOnAssetIDs = createStringSliceFromListInput(d.Get("exists_on_assets").(*schema.Set).List())
	}
	if d.HasChange("content") {
		content := createStringSliceFromListInput(d.Get("content").(*schema.Set).List())
		state.Content = &content
	}

	_, err = apiClient.ZeroTrust.UpdateState(*state)

	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceStateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*AuxoProvider)
	apiClient := provider.APIClient

	var diags diag.Diagnostics

	err := apiClient.ZeroTrust.DeleteStateByID(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
