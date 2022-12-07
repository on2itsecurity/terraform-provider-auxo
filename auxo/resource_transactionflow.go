package auxo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/on2itsecurity/go-auxo/zerotrust"
)

func resourceTransactionFlow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTransactionFlowCreate,
		ReadContext:   resourceTransactionFlowRead,
		UpdateContext: resourceTransactionFlowUpdate,
		DeleteContext: resourceTransactionFlowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"protectsurface": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The corresponding protect surface id",
			},
			"incoming_protectsurfaces_allow": {
				Type:        schema.TypeSet,
				Description: "The corresponding protect surfaces that are allowed to have incoming flows to this protect surface",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"incoming_protectsurfaces_block": {
				Type:        schema.TypeSet,
				Description: "The corresponding protect surfaces that are denied to have incoming flows to this protect surface",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"outgoing_protectsurfaces_allow": {
				Type:        schema.TypeSet,
				Description: "Protect surfaces to which flows are allowed (from this protect surface)",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"outgoing_protectsurfaces_block": {
				Type:        schema.TypeSet,
				Description: "Protect surfaces to which flows are denied (from this protect surface)",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceTransactionFlowCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	provider := m.(*AuxoProvider)
	apiClient := provider.APIClient

	psID := d.Get("protectsurface").(string)
	ps, _ := apiClient.ZeroTrust.GetProtectSurfaceByID(psID)

	//Flows from other PS
	flowsFromOtherPS := make(map[string]zerotrust.Flow)
	createFlowsFromSetInput(flowsFromOtherPS, d.Get("incoming_protectsurfaces_allow"), true)
	createFlowsFromSetInput(flowsFromOtherPS, d.Get("incoming_protectsurfaces_block"), false)
	ps.FlowsFromOtherPS = flowsFromOtherPS

	//Flows to other PS
	flowsToOtherPS := make(map[string]zerotrust.Flow)
	createFlowsFromSetInput(flowsToOtherPS, d.Get("outgoing_protectsurfaces_allow"), true)
	createFlowsFromSetInput(flowsToOtherPS, d.Get("outgoing_protectsurfaces_block"), false)
	ps.FlowsToOtherPS = flowsToOtherPS

	_, err := apiClient.ZeroTrust.CreateProtectSurfaceByObject(*ps, true)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(psID)

	resourceLocationRead(ctx, d, m)

	return diags
}

func resourceTransactionFlowRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	provider := m.(*AuxoProvider)
	apiClient := provider.APIClient

	ps, err := apiClient.ZeroTrust.GetProtectSurfaceByID(d.Get("protectsurface").(string))

	if err != nil {
		apiError := getAPIError(err)

		//NotExists
		if apiError.ID == "410" {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	d.SetId(ps.ID)

	d.Set("protectsurface", ps.ID)
	incomingPSAllow := []string{}
	incomingPSBlock := []string{}
	outgoingPSAllow := []string{}
	outgoingPSBlock := []string{}

	for psID, flow := range ps.FlowsFromOtherPS {
		if flow.Allow {
			incomingPSAllow = append(incomingPSAllow, psID)
		} else {
			incomingPSBlock = append(incomingPSBlock, psID)
		}
	}
	for psID, flow := range ps.FlowsToOtherPS {
		if flow.Allow {
			outgoingPSAllow = append(outgoingPSAllow, psID)
		} else {
			outgoingPSBlock = append(outgoingPSBlock, psID)
		}
	}

	d.Set("incoming_protectsurfaces_allow", incomingPSAllow)
	d.Set("incoming_protectsurfaces_block", incomingPSBlock)
	d.Set("outgoing_protectsurfaces_allow", outgoingPSAllow)
	d.Set("outgoing_protectsurfaces_block", outgoingPSBlock)

	return diags
}

func resourceTransactionFlowUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*AuxoProvider)
	apiClient := provider.APIClient

	ps, err := apiClient.ZeroTrust.GetProtectSurfaceByID(d.Get("protectsurface").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("incoming_protectsurfaces_allow") || d.HasChange("incoming_protectsurfaces_block") {
		flowsFromOtherPS := make(map[string]zerotrust.Flow)
		createFlowsFromSetInput(flowsFromOtherPS, d.Get("incoming_protectsurfaces_allow"), true)
		createFlowsFromSetInput(flowsFromOtherPS, d.Get("incoming_protectsurfaces_block"), false)
		ps.FlowsFromOtherPS = flowsFromOtherPS
	}
	if d.HasChange("outgoing_protectsurfaces_allow") || d.HasChange("outgoing_protectsurfaces_block") {
		flowsToOtherPS := make(map[string]zerotrust.Flow)
		createFlowsFromSetInput(flowsToOtherPS, d.Get("outgoing_protectsurfaces_allow"), true)
		createFlowsFromSetInput(flowsToOtherPS, d.Get("outgoing_protectsurfaces_block"), false)
		ps.FlowsToOtherPS = flowsToOtherPS
	}

	_, err = apiClient.ZeroTrust.CreateProtectSurfaceByObject(*ps, true)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceProtectSurfaceRead(ctx, d, m)
}

func resourceTransactionFlowDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*AuxoProvider)
	apiClient := provider.APIClient

	var diags diag.Diagnostics

	ps, err := apiClient.ZeroTrust.GetProtectSurfaceByID(d.Get("protectsurface").(string))

	ps.FlowsFromOtherPS = map[string]zerotrust.Flow{}
	ps.FlowsToOtherPS = map[string]zerotrust.Flow{}

	_, err = apiClient.ZeroTrust.CreateProtectSurfaceByObject(*ps, true)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

// createFlowsFromSetInput fills a map of flows from a set input
func createFlowsFromSetInput(flows map[string]zerotrust.Flow, inputField interface{}, allow bool) {
	for _, psID := range createStringSliceFromListInput(inputField.(*schema.Set).List()) {
		flows[psID] = zerotrust.Flow{Allow: allow}
	}
}
