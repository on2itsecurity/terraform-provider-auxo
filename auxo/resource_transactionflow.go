package auxo

import (
	"context"

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
	Protectsurface                 types.String   `tfsdk:"protectsurface"`
	Incoming_protectsurfaces_allow []types.String `tfsdk:"incoming_protectsurfaces_allow"`
	Incoming_protectsurfaces_block []types.String `tfsdk:"incoming_protectsurfaces_block"`
	Outgoing_protectsurfaces_allow []types.String `tfsdk:"outgoing_protectsurfaces_allow"`
	Outgoing_protectsurfaces_block []types.String `tfsdk:"outgoing_protectsurfaces_block"`
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

	for _, flow := range plan.Incoming_protectsurfaces_allow {
		ps.FlowsFromOtherPS[flow.ValueString()] = zerotrust.Flow{Allow: true}
	}

	for _, flow := range plan.Incoming_protectsurfaces_block {
		ps.FlowsFromOtherPS[flow.ValueString()] = zerotrust.Flow{Allow: false}
	}

	for _, flow := range plan.Outgoing_protectsurfaces_allow {
		ps.FlowsToOtherPS[flow.ValueString()] = zerotrust.Flow{Allow: true}
	}

	for _, flow := range plan.Outgoing_protectsurfaces_block {
		ps.FlowsToOtherPS[flow.ValueString()] = zerotrust.Flow{Allow: false}
	}

	ps, err = r.client.ZeroTrust.CreateProtectSurfaceByObject(*ps, true)

	if err != nil {
		resp.Diagnostics.AddError("Error creating transactionflow", "unexpected error: "+err.Error())
		return
	}

	// Map resonse to schema
	plan.Protectsurface = types.StringValue(ps.ID)

	incomingPSAllow := []basetypes.StringValue{}
	incomingPSBlock := []basetypes.StringValue{}
	outgoingPSAllow := []basetypes.StringValue{}
	outgoingPSBlock := []basetypes.StringValue{}

	for psID, flow := range ps.FlowsFromOtherPS {
		if flow.Allow {
			incomingPSAllow = append(incomingPSAllow, basetypes.NewStringValue(psID))
		} else {
			incomingPSBlock = append(incomingPSBlock, basetypes.NewStringValue(psID))
		}
	}

	for psID, flow := range ps.FlowsToOtherPS {
		if flow.Allow {
			outgoingPSAllow = append(outgoingPSAllow, basetypes.NewStringValue(psID))
		} else {
			outgoingPSBlock = append(outgoingPSBlock, basetypes.NewStringValue(psID))
		}
	}

	plan.Incoming_protectsurfaces_allow = incomingPSAllow
	plan.Incoming_protectsurfaces_block = incomingPSBlock
	plan.Outgoing_protectsurfaces_allow = outgoingPSAllow
	plan.Outgoing_protectsurfaces_block = outgoingPSBlock

	// Set state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *transactionflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *transactionflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *transactionflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
