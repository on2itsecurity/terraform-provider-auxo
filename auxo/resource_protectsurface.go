package auxo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/on2itsecurity/go-auxo"
	"github.com/on2itsecurity/go-auxo/zerotrust"
)

// var _ resource.ResourceWithModifyPlan = &protectsurfaceResource{}
var _ resource.Resource = &protectsurfaceResource{}

type protectsurfaceResource struct {
	client *auxo.Client
}

type protectsurfaceResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Uniqueness_key        types.String `tfsdk:"uniqueness_key"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	MainContact           types.String `tfsdk:"main_contact"`
	SecurityContact       types.String `tfsdk:"security_contact"`
	InControlBoundary     types.Bool   `tfsdk:"in_control_boundary"`
	InZeroTrustFocus      types.Bool   `tfsdk:"in_zero_trust_focus"`
	Relevance             types.Int64  `tfsdk:"relevance"`
	Confidentiality       types.Int64  `tfsdk:"confidentiality"`
	Integrity             types.Int64  `tfsdk:"integrity"`
	Availability          types.Int64  `tfsdk:"availability"`
	DataTags              types.Set    `tfsdk:"data_tags"`
	ComplianceTags        types.Set    `tfsdk:"compliance_tags"`
	CustomerLabels        types.Map    `tfsdk:"customer_labels"`
	SOCTags               types.Set    `tfsdk:"soc_tags"`
	AllowFlowsFromOutside types.Bool   `tfsdk:"allow_flows_from_outside"`
	AllowFlowsToOutside   types.Bool   `tfsdk:"allow_flows_to_outside"`
	MaturityStep1         types.Int64  `tfsdk:"maturity_step1"`
	MaturityStep2         types.Int64  `tfsdk:"maturity_step2"`
	MaturityStep3         types.Int64  `tfsdk:"maturity_step3"`
	MaturityStep4         types.Int64  `tfsdk:"maturity_step4"`
	MaturityStep5         types.Int64  `tfsdk:"maturity_step5"`
}

func NewProtectsurfaceResource() resource.Resource {
	return &protectsurfaceResource{}
}

func (r *protectsurfaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_protectsurface"
}

func (r *protectsurfaceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	// Retrieve the client from the provider config
	r.client = req.ProviderData.(*auxo.Client)
}

func (r *protectsurfaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "A zero trust protectsurface which reflects what you want to protect.",
		MarkdownDescription: "A zero trust protectsurface which reflects what you want to protect.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Computed unique ID of the resource protectsurface",
				MarkdownDescription: "Computed unique ID of the resource protectsurface",
				Required:            false,
				Optional:            false,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"uniqueness_key": schema.StringAttribute{
				Description:         "Custom and optinal uniqueness key to identify the resource protectsurface",
				MarkdownDescription: "Custom and optinal uniqueness key to identify the resource protectsurface",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the resource protectsurface",
				MarkdownDescription: "Name of the resource protectsurface",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the resource protectsurface",
				MarkdownDescription: "Description of the resource protectsurface",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"main_contact": schema.StringAttribute{
				Description:         "Main contact of the resource protectsurface",
				MarkdownDescription: "Main contact of the resource protectsurface",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"security_contact": schema.StringAttribute{
				Description:         "Security contact of the resource protectsurface",
				MarkdownDescription: "Security contact of the resource protectsurface",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"in_control_boundary": schema.BoolAttribute{
				Description:         "This protect surface is within the 'control boundary'",
				MarkdownDescription: "This protect surface is within the 'control boundary'",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"in_zero_trust_focus": schema.BoolAttribute{
				Description:         "This protect surface is within the 'zero trust focus' (actively maintained and monitored)",
				MarkdownDescription: "This protect surface is within the 'zero trust focus' (actively maintained and monitored)",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"relevance": schema.Int64Attribute{
				Description:         "Relevance of the resource protectsurface",
				MarkdownDescription: "Relevance of the resource protectsurface",
				Required:            true,
			},
			"confidentiality": schema.Int64Attribute{
				Description:         "Confidentiality of the resource protectsurface",
				MarkdownDescription: "Confidentiality of the resource protectsurface",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
			},
			"integrity": schema.Int64Attribute{
				Description:         "Integrity of the resource protectsurface",
				MarkdownDescription: "Integrity of the resource protectsurface",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
			},
			"availability": schema.Int64Attribute{
				Description:         "Availability of the resource protectsurface",
				MarkdownDescription: "Availability of the resource protectsurface",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
			},
			"data_tags": schema.SetAttribute{
				Description:         "Data tags of the resource protectsurface",
				MarkdownDescription: "Data tags of the resource protectsurface",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"compliance_tags": schema.SetAttribute{
				Description:         "Compliance tags of the resource protectsurface",
				MarkdownDescription: "Compliance tags of the resource protectsurface",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"customer_labels": schema.MapAttribute{
				Description:         "Customer labels of the resource protectsurface",
				MarkdownDescription: "Customer labels of the resource protectsurface",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"soc_tags": schema.SetAttribute{
				Description:         "SOC tags of the resource protectsurface, only use when advised by the SOC",
				MarkdownDescription: "SOC tags of the resource protectsurface, only use when advised by the SOC",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_flows_from_outside": schema.BoolAttribute{
				Description:         "Allow flows from outside of the protectsurface coming in",
				MarkdownDescription: "Allow flows from outside of the protectsurface coming in",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"allow_flows_to_outside": schema.BoolAttribute{
				Description:         "Allow flows to go outside of the protectsurface",
				MarkdownDescription: "Allow flows to go outside of the protectsurface",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"maturity_step1": schema.Int64Attribute{
				Description:         "Maturity step 1",
				MarkdownDescription: "Maturity step 1",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
			},
			"maturity_step2": schema.Int64Attribute{
				Description:         "Maturity step 2",
				MarkdownDescription: "Maturity step 2",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
			},
			"maturity_step3": schema.Int64Attribute{
				Description:         "Maturity step 3",
				MarkdownDescription: "Maturity step 3",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
			},
			"maturity_step4": schema.Int64Attribute{
				Description:         "Maturity step 4",
				MarkdownDescription: "Maturity step 4",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
			},
			"maturity_step5": schema.Int64Attribute{
				Description:         "Maturity step 5",
				MarkdownDescription: "Maturity step 5",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
			},
		},
	}
}

func (r *protectsurfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan protectsurfaceResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	protectsurface, d := resourceModelToProtectsurface(&plan, ctx, r)

	resp.Diagnostics.Append(d...)

	if resp.Diagnostics.HasError() {
		return
	}

	//Create the protectsurface
	result, err := r.client.ZeroTrust.CreateProtectSurfaceByObject(protectsurface, false)

	if err != nil {
		resp.Diagnostics.AddError("Error creating protect surface", "unexpected error: "+err.Error())
		return
	}

	//Map response to schema
	plan, _ = protectsurfaceToResourceModel(result, ctx)

	// Set state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *protectsurfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state protectsurfaceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed PS from AUXO
	result, err := r.client.ZeroTrust.GetProtectSurfaceByID(state.ID.ValueString())
	if err != nil {
		apiError := getAPIError(err)

		if apiError.ID == "410" { // Location not found and probably deleted
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("Error reading location", "unexpected error: "+err.Error())
			return
		}
	}

	//Overwrite state with refreshed PS
	state, _ = protectsurfaceToResourceModel(result, ctx)

	//Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *protectsurfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve values from plan
	var plan protectsurfaceResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	protectsurface, d := resourceModelToProtectsurface(&plan, ctx, r)

	resp.Diagnostics.Append(d...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.ZeroTrust.UpdateProtectSurface(protectsurface)

	if err != nil {
		resp.Diagnostics.AddError("Error updating protect surface", "unexpected error: "+err.Error())
		return
	}

	plan, _ = protectsurfaceToResourceModel(result, ctx)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *protectsurfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve values from state
	var ps protectsurfaceResourceModel

	diags := req.State.Get(ctx, &ps)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.ZeroTrust.DeleteProtectSurfaceByID(ps.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Error deleting protect surface", "unexpected error: "+err.Error())
		return
	}

}

// resourceModelToProtectsurface maps the resource model to the zerotrust.protectsurface object
func resourceModelToProtectsurface(plan *protectsurfaceResourceModel, ctx context.Context, r *protectsurfaceResource) (zerotrust.ProtectSurface, diag.Diagnostics) {
	var diag diag.Diagnostics
	var st, dt, ct []string

	if !plan.DataTags.IsNull() {
		_ = plan.DataTags.ElementsAs(ctx, &dt, false)
	}
	if !plan.ComplianceTags.IsNull() {
		_ = plan.ComplianceTags.ElementsAs(ctx, &ct, false)
	}
	if !plan.SOCTags.IsNull() {
		_ = plan.SOCTags.ElementsAs(ctx, &st, false)
	}
	var cl map[string]string
	types.Map.ElementsAs(plan.CustomerLabels, ctx, &cl, false)

	//Create the protectsurface object
	protectsurface := zerotrust.ProtectSurface{
		ID:                      plan.ID.ValueString(),
		UniquenessKey:           plan.Uniqueness_key.ValueString(),
		Name:                    plan.Name.ValueString(),
		Description:             plan.Description.ValueString(),
		MainContactPersonID:     plan.MainContact.ValueString(),
		SecurityContactPersonID: plan.SecurityContact.ValueString(),
		InControlBoundary:       plan.InControlBoundary.ValueBool(),
		InZeroTrustFocus:        plan.InZeroTrustFocus.ValueBool(),
		Relevance:               int(plan.Relevance.ValueInt64()),
		Confidentiality:         int(plan.Confidentiality.ValueInt64()),
		Integrity:               int(plan.Integrity.ValueInt64()),
		Availability:            int(plan.Availability.ValueInt64()),
		DataTags:                dt,
		ComplianceTags:          ct,
		CustomerLabels:          cl,
		SocTags:                 st,
		FlowsFromOutside: zerotrust.Flow{
			Allow: plan.AllowFlowsFromOutside.ValueBool(),
		},
		FlowsToOutside: zerotrust.Flow{
			Allow: plan.AllowFlowsToOutside.ValueBool(),
		},
		Maturity: zerotrust.Maturity{
			Step1: int(plan.MaturityStep1.ValueInt64()),
			Step2: int(plan.MaturityStep2.ValueInt64()),
			Step3: int(plan.MaturityStep3.ValueInt64()),
			Step4: int(plan.MaturityStep4.ValueInt64()),
			Step5: int(plan.MaturityStep5.ValueInt64()),
		},
	}

	return protectsurface, diag
}

// protectsurfaceToResourceModel maps the zerotrust.protectsurface object to the resource model
func protectsurfaceToResourceModel(ps *zerotrust.ProtectSurface, ctx context.Context) (protectsurfaceResourceModel, diag.Diagnostics) {
	cl, diag := types.MapValueFrom(ctx, types.StringType, ps.CustomerLabels)

	var st, dt, ct basetypes.SetValue
	if ps.ComplianceTags != nil {
		ct, _ = types.SetValueFrom(ctx, types.StringType, ps.ComplianceTags)
	}
	if ps.DataTags != nil {
		dt, _ = types.SetValueFrom(ctx, types.StringType, ps.DataTags)
	}
	if ps.SocTags != nil {
		st, _ = types.SetValueFrom(ctx, types.StringType, ps.SocTags)
	}

	psrm := protectsurfaceResourceModel{
		ID:                    types.StringValue(ps.ID),
		Uniqueness_key:        types.StringValue(ps.UniquenessKey),
		Name:                  types.StringValue(ps.Name),
		Description:           types.StringValue(ps.Description),
		MainContact:           types.StringValue(ps.MainContactPersonID),
		SecurityContact:       types.StringValue(ps.SecurityContactPersonID),
		InControlBoundary:     types.BoolValue(ps.InControlBoundary),
		InZeroTrustFocus:      types.BoolValue(ps.InZeroTrustFocus),
		Relevance:             types.Int64Value(int64(ps.Relevance)),
		Confidentiality:       types.Int64Value(int64(ps.Confidentiality)),
		Integrity:             types.Int64Value(int64(ps.Integrity)),
		Availability:          types.Int64Value(int64(ps.Availability)),
		DataTags:              dt,
		ComplianceTags:        ct,
		CustomerLabels:        cl,
		SOCTags:               st,
		AllowFlowsFromOutside: types.BoolValue(ps.FlowsFromOutside.Allow),
		AllowFlowsToOutside:   types.BoolValue(ps.FlowsToOutside.Allow),
		MaturityStep1:         types.Int64Value(int64(ps.Maturity.Step1)),
		MaturityStep2:         types.Int64Value(int64(ps.Maturity.Step2)),
		MaturityStep3:         types.Int64Value(int64(ps.Maturity.Step3)),
		MaturityStep4:         types.Int64Value(int64(ps.Maturity.Step4)),
		MaturityStep5:         types.Int64Value(int64(ps.Maturity.Step5)),
	}

	return psrm, diag
}
