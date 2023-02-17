package auxo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/on2itsecurity/go-auxo"
	"github.com/on2itsecurity/go-auxo/zerotrust"
)

var _ resource.Resource = &protectsurfaceResource{}

type protectsurfaceResource struct {
	client *auxo.Client
}

type protectsurfaceResourceModel struct {
	ID                types.String   `tfsdk:"id"`
	Uniqueness_key    types.String   `tfsdk:"uniqueness_key"`
	Name              types.String   `tfsdk:"name"`
	Description       types.String   `tfsdk:"description"`
	MainContact       types.String   `tfsdk:"main_contact"`
	SecurityContact   types.String   `tfsdk:"security_contact"`
	InControlBoundary types.Bool     `tfsdk:"in_control_boundary"`
	InZeroTrustFocus  types.Bool     `tfsdk:"in_zero_trust_focus"`
	Relevance         types.Int64    `tfsdk:"relevance"`
	Confidentiality   types.Int64    `tfsdk:"confidentiality"`
	Integrity         types.Int64    `tfsdk:"integrity"`
	Availability      types.Int64    `tfsdk:"availability"`
	DataTags          []types.String `tfsdk:"data_tags"`
	ComplianceTags    []types.String `tfsdk:"compliance_tags"`
	//CustomerLabels        types.Map      `tfsdk:"customer_labels"`
	CustomerLabels        map[string]string `tfsdk:"customer_labels"`
	SOCTags               []types.String    `tfsdk:"soc_tags"`
	AllowFlowsFromOutside types.Bool        `tfsdk:"allow_flows_from_outside"`
	AllowFlowsToOutside   types.Bool        `tfsdk:"allow_flows_to_outside"`
	MaturityStep1         types.Int64       `tfsdk:"maturity_step1"`
	MaturityStep2         types.Int64       `tfsdk:"maturity_step2"`
	MaturityStep3         types.Int64       `tfsdk:"maturity_step3"`
	MaturityStep4         types.Int64       `tfsdk:"maturity_step4"`
	MaturityStep5         types.Int64       `tfsdk:"maturity_step5"`
	//TODO Measures
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
		Description:         "Auxo Protectsurface",
		MarkdownDescription: "Auxo Protectsurface",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Computed unique ID of the resource protectsurface",
				MarkdownDescription: "Computed unique ID of the resource protectsurface",
				Computed:            true,
			},
			"uniqueness_key": schema.StringAttribute{
				Description:         "Custom and optinal uniqueness key to identify the resource protectsurface",
				MarkdownDescription: "Custom and optinal uniqueness key to identify the resource protectsurface",
				Optional:            true,
				Computed:            true,
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
			},
			"main_contact": schema.StringAttribute{
				Description:         "Main contact of the resource protectsurface",
				MarkdownDescription: "Main contact of the resource protectsurface",
				Optional:            true,
			},
			"security_contact": schema.StringAttribute{
				Description:         "Security contact of the resource protectsurface",
				MarkdownDescription: "Security contact of the resource protectsurface",
				Optional:            true,
			},
			"in_control_boundary": schema.BoolAttribute{
				Description:         "This protect surface is within the 'control boundary'",
				MarkdownDescription: "This protect surface is within the 'control boundary'",
				Optional:            true,
			},
			"in_zero_trust_focus": schema.BoolAttribute{
				Description:         "This protect surface is within the 'zero trust focus' (actively maintained and monitored",
				MarkdownDescription: "This protect surface is within the 'zero trust focus' (actively maintained and monitored",
				Optional:            true,
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
			},
			"integrity": schema.Int64Attribute{
				Description:         "Integrity of the resource protectsurface",
				MarkdownDescription: "Integrity of the resource protectsurface",
				Optional:            true,
			},
			"availability": schema.Int64Attribute{
				Description:         "Availability of the resource protectsurface",
				MarkdownDescription: "Availability of the resource protectsurface",
				Optional:            true,
			},
			"data_tags": schema.SetAttribute{
				Description:         "Data tags of the resource protectsurface",
				MarkdownDescription: "Data tags of the resource protectsurface",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"compliance_tags": schema.SetAttribute{
				Description:         "Compliance tags of the resource protectsurface",
				MarkdownDescription: "Compliance tags of the resource protectsurface",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"customer_labels": schema.MapAttribute{
				Description:         "Customer labels of the resource protectsurface",
				MarkdownDescription: "Customer labels of the resource protectsurface",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"soc_tags": schema.SetAttribute{
				Description:         "Soc tags of the resource protectsurface",
				MarkdownDescription: "Soc tags of the resource protectsurface",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"allow_flows_from_outside": schema.BoolAttribute{
				Description:         "Allow flows from outside of the protectsurface coming in",
				MarkdownDescription: "Allow flows from outside of the protectsurface coming in",
				Optional:            true,
			},
			"allow_flows_to_outside": schema.BoolAttribute{
				Description:         "Allow flows to go outside of the protectsurface",
				MarkdownDescription: "Allow flows to go outside of the protectsurface",
				Optional:            true,
			},
			"maturity_step1": schema.Int64Attribute{
				Description:         "Maturity step 1",
				MarkdownDescription: "Maturity step 1",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					Int64DefaultValue(types.Int64Value(1)),
				},
			},
			"maturity_step2": schema.Int64Attribute{
				Description:         "Maturity step 2",
				MarkdownDescription: "Maturity step 2",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					Int64DefaultValue(types.Int64Value(1)),
				},
			},
			"maturity_step3": schema.Int64Attribute{
				Description:         "Maturity step 3",
				MarkdownDescription: "Maturity step 3",
				Optional:            true,
				Computed:            true,
			},
			"maturity_step4": schema.Int64Attribute{
				Description:         "Maturity step 4",
				MarkdownDescription: "Maturity step 4",
				Optional:            true,
				Computed:            true,
			},
			"maturity_step5": schema.Int64Attribute{
				Description:         "Maturity step 5",
				MarkdownDescription: "Maturity step 5",
				Optional:            true,
				Computed:            true,
			},
			//TODO Measures
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

	//Get slices from plan
	dt := getSliceFromSetOfString(plan.DataTags)
	ct := getSliceFromSetOfString(plan.ComplianceTags)
	st := getSliceFromSetOfString(plan.SOCTags)

	// cl := make(map[string]string)
	// for k, v := range plan.CustomerLabels {
	// 	cl[k] = v.ValueString()
	// }

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
		CustomerLabels:          plan.CustomerLabels,
		SocTags:                 st,
		FlowsFromOutside: zerotrust.Flow{
			Allow: plan.AllowFlowsFromOutside.ValueBool(),
		},
		FlowsToOutside: zerotrust.Flow{
			Allow: plan.AllowFlowsToOutside.ValueBool(),
		},
		Maturity: zerotrust.Maturity{
			Step1: int(plan.MaturityStep2.ValueInt64()),
			Step2: int(plan.MaturityStep2.ValueInt64()),
			Step3: int(plan.MaturityStep3.ValueInt64()),
			Step4: int(plan.MaturityStep4.ValueInt64()),
			Step5: int(plan.MaturityStep5.ValueInt64()),
		},
		//TODO Measures
	}

	//Create the protectsurface
	result, err := r.client.ZeroTrust.CreateProtectSurfaceByObject(protectsurface, false)

	if err != nil {
		resp.Diagnostics.AddError("Error creating protect surface", "unexpected error: "+err.Error())
	}

	//Map response to schema
	plan.ID = types.StringValue(result.ID)
	plan.Uniqueness_key = types.StringValue(result.UniquenessKey)
	plan.Name = types.StringValue(result.Name)
	plan.Description = types.StringValue(result.Description)
	plan.MainContact = types.StringValue(result.MainContactPersonID)
	plan.SecurityContact = types.StringValue(result.SecurityContactPersonID)
	plan.InControlBoundary = types.BoolValue(result.InControlBoundary)
	plan.InZeroTrustFocus = types.BoolValue(result.InZeroTrustFocus)
	plan.Relevance = types.Int64Value(int64(result.Relevance))
	plan.Confidentiality = types.Int64Value(int64(result.Confidentiality))
	plan.Integrity = types.Int64Value(int64(result.Integrity))
	plan.Availability = types.Int64Value(int64(result.Availability))
	plan.DataTags = getSetOfStringFromSlice(result.DataTags)
	plan.ComplianceTags = getSetOfStringFromSlice(result.ComplianceTags)
	//plan.CustomerLabels = result.CustomerLabels
	plan.SOCTags = getSetOfStringFromSlice(result.SocTags)
	plan.AllowFlowsFromOutside = types.BoolValue(result.FlowsFromOutside.Allow)
	plan.AllowFlowsToOutside = types.BoolValue(result.FlowsToOutside.Allow)
	plan.MaturityStep1 = types.Int64Value(int64(result.Maturity.Step1))
	plan.MaturityStep2 = types.Int64Value(int64(result.Maturity.Step2))
	plan.MaturityStep3 = types.Int64Value(int64(result.Maturity.Step3))
	plan.MaturityStep4 = types.Int64Value(int64(result.Maturity.Step4))
	plan.MaturityStep5 = types.Int64Value(int64(result.Maturity.Step5))
	//TODO Measures

	// Set state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *protectsurfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *protectsurfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *protectsurfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
