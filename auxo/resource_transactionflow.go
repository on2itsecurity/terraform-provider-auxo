package auxo

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/on2itsecurity/go-auxo"
	"github.com/on2itsecurity/go-auxo/zerotrust"
)

var _ resource.Resource = &transactionflowResource{}

type transactionflowResource struct {
	client *auxo.Client
}

type transactionflowResourceModel struct {
	Protectsurface                 types.String `tfsdk:"protectsurface"`
	Incoming_protectsurfaces_allow types.Set    `tfsdk:"incoming_protectsurfaces_allow"`
	Incoming_protectsurfaces_block types.Set    `tfsdk:"incoming_protectsurfaces_block"`
	Outgoing_protectsurfaces_allow types.Set    `tfsdk:"outgoing_protectsurfaces_allow"`
	Outgoing_protectsurfaces_block types.Set    `tfsdk:"outgoing_protectsurfaces_block"`
}

type flows struct {
	incomingPSAllow []basetypes.StringValue
	incomingPSBlock []basetypes.StringValue
	outgoingPSAllow []basetypes.StringValue
	outgoingPSBlock []basetypes.StringValue
}

func NewTransactionflowResource() resource.Resource {
	return &transactionflowResource{}
}

func (r *transactionflowResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_transactionflow"
}

func (r *transactionflowResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	// Retrieve the client from the provider config
	r.client = req.ProviderData.(*auxo.Client)
}

func (r *transactionflowResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Auxo transactionflows",
		MarkdownDescription: "Auxo transactionflows",
		Attributes: map[string]schema.Attribute{
			"protectsurface": schema.StringAttribute{
				Description:         "The ID of the protectsurface",
				MarkdownDescription: "The ID of the protectsurface",
				Required:            true,
			},
			"incoming_protectsurfaces_allow": schema.SetAttribute{
				Description:         "The IDs of the protectsurface that are allowed to send traffic to this protectsurface",
				MarkdownDescription: "The IDs of the protectsurface that are allowed to send traffic to this protectsurface",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"incoming_protectsurfaces_block": schema.SetAttribute{
				Description:         "The IDs of the protectsurface that are blocked to send traffic to this protectsurface",
				MarkdownDescription: "The IDs of the protectsurface that are blocked to send traffic to this protectsurface",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"outgoing_protectsurfaces_allow": schema.SetAttribute{
				Description:         "The IDs of the protectsurface that are allowed to send traffic to from this protectsurface",
				MarkdownDescription: "The IDs of the protectsurface that are allowed to send traffic to from this protectsurface",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"outgoing_protectsurfaces_block": schema.SetAttribute{
				Description:         "The IDs of the protectsurface that are blocked to send traffic to from this protectsurface",
				MarkdownDescription: "The IDs of the protectsurface that are blocked to send traffic to from this protectsurface",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *transactionflowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan transactionflowResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// create transactionflow
	psID := plan.Protectsurface.ValueString()
	ps, err := r.client.ZeroTrust.GetProtectSurfaceByID(psID)

	if err != nil {
		resp.Diagnostics.AddError("Error creating transactionflow", "unexpected error: "+err.Error())
		return
	}

	var f flows
	_ = plan.Incoming_protectsurfaces_allow.ElementsAs(ctx, &f.incomingPSAllow, false)
	_ = plan.Incoming_protectsurfaces_block.ElementsAs(ctx, &f.incomingPSBlock, false)
	_ = plan.Outgoing_protectsurfaces_allow.ElementsAs(ctx, &f.outgoingPSAllow, false)
	_ = plan.Outgoing_protectsurfaces_block.ElementsAs(ctx, &f.outgoingPSBlock, false)

	ps, err = setFlowsOnPS(ps, f)
	if err != nil {
		resp.Diagnostics.AddError("Error creating transactionflow", err.Error())
		return
	}

	ps, err = r.client.ZeroTrust.CreateProtectSurfaceByObject(*ps, true)

	if err != nil {
		resp.Diagnostics.AddError("Error creating transactionflow", "unexpected error: "+err.Error())
		return
	}

	// Map resonse to schema
	plan.Protectsurface = types.StringValue(ps.ID)
	f = readFlowsFromPS(ps)
	plan.Incoming_protectsurfaces_allow, _ = types.SetValueFrom(ctx, types.StringType, f.incomingPSAllow)
	plan.Incoming_protectsurfaces_block, _ = types.SetValueFrom(ctx, types.StringType, f.incomingPSBlock)
	plan.Outgoing_protectsurfaces_allow, _ = types.SetValueFrom(ctx, types.StringType, f.outgoingPSAllow)
	plan.Outgoing_protectsurfaces_block, _ = types.SetValueFrom(ctx, types.StringType, f.outgoingPSBlock)

	// Set state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *transactionflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state transactionflowResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed state from AUXO
	result, err := r.client.ZeroTrust.GetProtectSurfaceByID(state.Protectsurface.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading location", "unexpected error: "+err.Error())
		return
	}

	//Overwrite state with refreshed state
	f := readFlowsFromPS(result)
	state.Protectsurface = types.StringValue(result.ID)
	state.Incoming_protectsurfaces_allow, _ = types.SetValueFrom(ctx, types.StringType, f.incomingPSAllow)
	state.Incoming_protectsurfaces_block, _ = types.SetValueFrom(ctx, types.StringType, f.incomingPSBlock)
	state.Outgoing_protectsurfaces_allow, _ = types.SetValueFrom(ctx, types.StringType, f.outgoingPSAllow)
	state.Outgoing_protectsurfaces_block, _ = types.SetValueFrom(ctx, types.StringType, f.outgoingPSBlock)

	//Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *transactionflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Retrieve values from plan
	var plan transactionflowResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// create transactionflow
	psID := plan.Protectsurface.ValueString()
	ps, err := r.client.ZeroTrust.GetProtectSurfaceByID(psID)

	if err != nil {
		resp.Diagnostics.AddError("Error creating transactionflow", "unexpected error: "+err.Error())
		return
	}

	var f flows
	_ = plan.Incoming_protectsurfaces_allow.ElementsAs(ctx, &f.incomingPSAllow, false)
	_ = plan.Incoming_protectsurfaces_block.ElementsAs(ctx, &f.incomingPSBlock, false)
	_ = plan.Outgoing_protectsurfaces_allow.ElementsAs(ctx, &f.outgoingPSAllow, false)
	_ = plan.Outgoing_protectsurfaces_block.ElementsAs(ctx, &f.outgoingPSBlock, false)

	ps, err = setFlowsOnPS(ps, f)
	if err != nil {
		resp.Diagnostics.AddError("Error creating transactionflow", err.Error())
		return
	}

	ps, err = r.client.ZeroTrust.CreateProtectSurfaceByObject(*ps, true)

	if err != nil {
		resp.Diagnostics.AddError("Error creating transactionflow", "unexpected error: "+err.Error())
		return
	}

	// Map resonse to schema
	plan.Protectsurface = types.StringValue(ps.ID)
	f = readFlowsFromPS(ps)

	plan.Incoming_protectsurfaces_allow, _ = types.SetValueFrom(ctx, types.StringType, f.incomingPSAllow)
	plan.Incoming_protectsurfaces_block, _ = types.SetValueFrom(ctx, types.StringType, f.incomingPSBlock)
	plan.Outgoing_protectsurfaces_allow, _ = types.SetValueFrom(ctx, types.StringType, f.outgoingPSAllow)
	plan.Outgoing_protectsurfaces_block, _ = types.SetValueFrom(ctx, types.StringType, f.outgoingPSBlock)

	// Set state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *transactionflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state transactionflowResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get PS and remove flows
	ps, err := r.client.ZeroTrust.GetProtectSurfaceByID(state.Protectsurface.ValueString())
	ps.FlowsFromOtherPS = map[string]zerotrust.Flow{}
	ps.FlowsToOtherPS = map[string]zerotrust.Flow{}

	// Update PS, with deleted flows
	_, err = r.client.ZeroTrust.CreateProtectSurfaceByObject(*ps, true)

	if err != nil {
		resp.Diagnostics.AddError("Error deleting state", "unexpected error: "+err.Error())
		return
	}
}

// readFlowsFromPS, get a ProtectSurface and return a flows struct, which can be used to map directly on plan & state
func readFlowsFromPS(ps *zerotrust.ProtectSurface) flows {
	var f flows

	for psID, flow := range ps.FlowsFromOtherPS {
		if flow.Allow {
			f.incomingPSAllow = append(f.incomingPSAllow, basetypes.NewStringValue(psID))
		} else {
			f.incomingPSBlock = append(f.incomingPSBlock, basetypes.NewStringValue(psID))
		}
	}

	for psID, flow := range ps.FlowsToOtherPS {
		if flow.Allow {
			f.outgoingPSAllow = append(f.outgoingPSAllow, basetypes.NewStringValue(psID))
		} else {
			f.outgoingPSBlock = append(f.outgoingPSBlock, basetypes.NewStringValue(psID))
		}
	}

	return f
}

// setFlowsOnPS, get a ProtectSurface and a flows struct, and set the flows on the ProtectSurface
func setFlowsOnPS(ps *zerotrust.ProtectSurface, flows flows) (*zerotrust.ProtectSurface, error) {

	//Check for duplicates in incoming and outgoing, which is not allowed
	for _, flow := range flows.incomingPSAllow {
		if sliceContains(getSliceFromSetOfString(flows.incomingPSBlock), flow.ValueString()) {
			return nil, fmt.Errorf("duplicate in incoming_protectsurfaces_allow and incoming_protectsurfaces_block, protectsurface ID: %s", flow.ValueString())
		}
	}

	for _, flow := range flows.outgoingPSAllow {
		if sliceContains(getSliceFromSetOfString(flows.outgoingPSBlock), flow.ValueString()) {
			return nil, fmt.Errorf("duplicate in outgoing_protectsurfaces_allow and outgoing_protectsurfaces_block, protectsurface ID: %s", flow.ValueString())
		}
	}

	ps.FlowsFromOtherPS = map[string]zerotrust.Flow{}
	ps.FlowsToOtherPS = map[string]zerotrust.Flow{}

	for _, flow := range flows.incomingPSAllow {
		ps.FlowsFromOtherPS[flow.ValueString()] = zerotrust.Flow{Allow: true}
	}

	for _, flow := range flows.incomingPSBlock {
		ps.FlowsFromOtherPS[flow.ValueString()] = zerotrust.Flow{Allow: false}
	}

	for _, flow := range flows.outgoingPSAllow {
		ps.FlowsToOtherPS[flow.ValueString()] = zerotrust.Flow{Allow: true}
	}

	for _, flow := range flows.outgoingPSBlock {
		ps.FlowsToOtherPS[flow.ValueString()] = zerotrust.Flow{Allow: false}
	}

	return ps, nil
}
