package auxo

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/on2itsecurity/go-auxo"
	"github.com/on2itsecurity/go-auxo/zerotrust"
)

var _ resource.Resource = &protectsurfaceResource{}

type protectsurfaceResource struct {
	client *auxo.Client
}

type protectsurfaceResourceModel struct {
	ID                    types.String       `tfsdk:"id"`
	Uniqueness_key        types.String       `tfsdk:"uniqueness_key"`
	Name                  types.String       `tfsdk:"name"`
	Description           types.String       `tfsdk:"description"`
	MainContact           types.String       `tfsdk:"main_contact"`
	SecurityContact       types.String       `tfsdk:"security_contact"`
	InControlBoundary     types.Bool         `tfsdk:"in_control_boundary"`
	InZeroTrustFocus      types.Bool         `tfsdk:"in_zero_trust_focus"`
	Relevance             types.Int64        `tfsdk:"relevance"`
	Confidentiality       types.Int64        `tfsdk:"confidentiality"`
	Integrity             types.Int64        `tfsdk:"integrity"`
	Availability          types.Int64        `tfsdk:"availability"`
	DataTags              []types.String     `tfsdk:"data_tags"`
	ComplianceTags        []types.String     `tfsdk:"compliance_tags"`
	CustomerLabels        types.Map          `tfsdk:"customer_labels"`
	SOCTags               []types.String     `tfsdk:"soc_tags"`
	AllowFlowsFromOutside types.Bool         `tfsdk:"allow_flows_from_outside"`
	AllowFlowsToOutside   types.Bool         `tfsdk:"allow_flows_to_outside"`
	MaturityStep1         types.Int64        `tfsdk:"maturity_step1"`
	MaturityStep2         types.Int64        `tfsdk:"maturity_step2"`
	MaturityStep3         types.Int64        `tfsdk:"maturity_step3"`
	MaturityStep4         types.Int64        `tfsdk:"maturity_step4"`
	MaturityStep5         types.Int64        `tfsdk:"maturity_step5"`
	Measures              map[string]measure `tfsdk:"measures"`
}

type measure struct {
	//Measure               types.String `tfsdk:"measure"`
	Assigned              types.Bool   `tfsdk:"assigned"`
	Assigned_by           types.String `tfsdk:"assigned_by"`
	Assigned_timestamp    types.Int64  `tfsdk:"assigned_timestamp"`
	Implemented           types.Bool   `tfsdk:"implemented"`
	Implemented_by        types.String `tfsdk:"implemented_by"`
	Implemented_timestamp types.Int64  `tfsdk:"implemented_timestamp"`
	Evidenced             types.Bool   `tfsdk:"evidenced"`
	Evidenced_by          types.String `tfsdk:"evidenced_by"`
	Evidenced_timestamp   types.Int64  `tfsdk:"evidenced_timestamp"`
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
				Description:         "This protect surface is within the 'zero trust focus' (actively maintained and monitored",
				MarkdownDescription: "This protect surface is within the 'zero trust focus' (actively maintained and monitored",
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
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
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
			"measures": schema.MapNestedAttribute{
				Description:         "Measures of the resource protectsurface",
				MarkdownDescription: "Measures of the resource protectsurface",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"assigned": schema.BoolAttribute{
							Description:         "Measure assigned to the protectsurface",
							MarkdownDescription: "Measure assigned to the protectsurface",
							Required:            true,
						},
						"assigned_by": schema.StringAttribute{
							Description:         "Who assigned this measure to the protectsurface",
							MarkdownDescription: "Who assigned this measure to the protectsurface",
							Optional:            true,
							Computed:            true,
							Default:             stringdefault.StaticString(""),
						},
						"assigned_timestamp": schema.Int64Attribute{
							Description:         "When was this measure assigned to the protectsurface",
							MarkdownDescription: "When was this measure assigned to the protectsurface",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"implemented": schema.BoolAttribute{
							Description:         "Is this measure implemented to the protectsurface",
							MarkdownDescription: "Is this measure implemented to the protectsurface",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
						"implemented_by": schema.StringAttribute{
							Description:         "Who implemented this measure to the protectsurface",
							MarkdownDescription: "Who implemented this measure to the protectsurface",
							Optional:            true,
							Computed:            true,
							Default:             stringdefault.StaticString(""),
						},
						"implemented_timestamp": schema.Int64Attribute{
							Description:         "When was this measure implemented to the protectsurface",
							MarkdownDescription: "When was this measure implemented to the protectsurface",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"evidenced": schema.BoolAttribute{
							Description:         "Is there evidence that this measure is implemented",
							MarkdownDescription: "Is there evidence that this measure is implemented",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
						"evidenced_by": schema.StringAttribute{
							Description:         "Who evidenced that this measure is implementd",
							MarkdownDescription: "Who evidenced that this measure is implementd",
							Optional:            true,
							Computed:            true,
							Default:             stringdefault.StaticString(""),
						},
						"evidenced_timestamp": schema.Int64Attribute{
							Description:         "When was this measure evidenced",
							MarkdownDescription: "When was this measure evidenced",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
					},
				},
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

	//Get slices from plan
	dt := getSliceFromSetOfString(plan.DataTags)
	ct := getSliceFromSetOfString(plan.ComplianceTags)
	st := getSliceFromSetOfString(plan.SOCTags)
	//Get map from plan
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

	//Measures
	measureMap := make(map[string]zerotrust.MeasureState, 0)

	//Loop through measures
	for k, m := range plan.Measures {

		//Check if measure exists
		if !sliceContains(r.getAvailableMeasures(), k) {
			resp.Diagnostics.AddError("Measure does not exists.",
				"Messure ["+k+"] does not exist, available measures ["+strings.Join(r.getAvailableMeasures(), ",")+"]")
			return
		}

		var assigned_timestamp int
		if !(m.Assigned_timestamp.IsUnknown() || m.Assigned_timestamp.IsNull()) {
			assigned_timestamp = int(m.Assigned_timestamp.ValueInt64())
		} else {
			assigned_timestamp = int(time.Now().Unix())
		}

		assignment := zerotrust.Assignment{
			Assigned:                 m.Assigned.ValueBool(),
			LastDeterminedByPersonID: m.Assigned_by.ValueString(),
			LastDeterminedTimestamp:  assigned_timestamp,
		}

		var implemented_timestamp int
		if !(m.Implemented_timestamp.IsUnknown() || m.Implemented_timestamp.IsNull()) {
			implemented_timestamp = int(m.Implemented_timestamp.ValueInt64())
		} else {
			implemented_timestamp = int(time.Now().Unix())
		}

		implementation := zerotrust.Implementation{
			Implemented:              m.Implemented.ValueBool(),
			LastDeterminedByPersonID: m.Implemented_by.ValueString(),
			LastDeterminedTimestamp:  implemented_timestamp,
		}

		var evidenced_timestamp int
		if !(m.Evidenced_timestamp.IsUnknown() || m.Evidenced_timestamp.IsNull()) {
			evidenced_timestamp = int(m.Evidenced_timestamp.ValueInt64())
		} else {
			evidenced_timestamp = int(time.Now().Unix())
		}

		evidence := zerotrust.Evidence{
			Evidenced:                m.Evidenced.ValueBool(),
			LastDeterminedByPersonID: m.Evidenced_by.ValueString(),
			LastDeterminedTimestamp:  evidenced_timestamp,
		}

		measureMap[k] = zerotrust.MeasureState{
			Assignment:     &assignment,
			Implementation: &implementation,
			Evidence:       &evidence,
		}
	}
	if len(measureMap) == 0 {
		measureMap = nil
	}
	protectsurface.Measures = measureMap

	//Create the protectsurface
	result, err := r.client.ZeroTrust.CreateProtectSurfaceByObject(protectsurface, false)

	// //DEBUG
	// var prettyJSON bytes.Buffer
	// byteOutput, _ := json.Marshal(protectsurface)
	// json.Indent(&prettyJSON, byteOutput, "", "\t")
	// resp.Diagnostics.AddError("PS", prettyJSON.String())
	// //DEBUG

	if err != nil {
		resp.Diagnostics.AddError("Error creating protect surface", "unexpected error: "+err.Error())
		return
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
	plan.CustomerLabels, _ = types.MapValueFrom(ctx, types.StringType, result.CustomerLabels)
	plan.SOCTags = getSetOfStringFromSlice(result.SocTags)
	plan.AllowFlowsFromOutside = types.BoolValue(result.FlowsFromOutside.Allow)
	plan.AllowFlowsToOutside = types.BoolValue(result.FlowsToOutside.Allow)
	plan.MaturityStep1 = types.Int64Value(int64(result.Maturity.Step1))
	plan.MaturityStep2 = types.Int64Value(int64(result.Maturity.Step2))
	plan.MaturityStep3 = types.Int64Value(int64(result.Maturity.Step3))
	plan.MaturityStep4 = types.Int64Value(int64(result.Maturity.Step4))
	plan.MaturityStep5 = types.Int64Value(int64(result.Maturity.Step5))
	plan.Measures = getMeasuresFromMap(result.Measures)

	// Set state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *protectsurfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var ps protectsurfaceResourceModel
	diags := req.State.Get(ctx, &ps)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed PS from AUXO
	result, err := r.client.ZeroTrust.GetProtectSurfaceByID(ps.ID.ValueString())
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
	//Get map from state
	//ID and UK cannot have changed
	ps.Name = types.StringValue(result.Name)
	ps.Description = types.StringValue(result.Description)
	ps.MainContact = types.StringValue(result.MainContactPersonID)
	ps.SecurityContact = types.StringValue(result.SecurityContactPersonID)
	ps.InControlBoundary = types.BoolValue(result.InControlBoundary)
	ps.InZeroTrustFocus = types.BoolValue(result.InZeroTrustFocus)
	ps.Relevance = types.Int64Value(int64(result.Relevance))
	ps.Confidentiality = types.Int64Value(int64(result.Confidentiality))
	ps.Integrity = types.Int64Value(int64(result.Integrity))
	ps.Availability = types.Int64Value(int64(result.Availability))
	ps.DataTags = getSetOfStringFromSlice(result.DataTags)
	ps.ComplianceTags = getSetOfStringFromSlice(result.ComplianceTags)
	ps.CustomerLabels, _ = types.MapValueFrom(ctx, types.StringType, result.CustomerLabels)
	ps.SOCTags = getSetOfStringFromSlice(result.SocTags)
	ps.AllowFlowsFromOutside = types.BoolValue(result.FlowsFromOutside.Allow)
	ps.AllowFlowsToOutside = types.BoolValue(result.FlowsToOutside.Allow)
	ps.MaturityStep1 = types.Int64Value(int64(result.Maturity.Step1))
	ps.MaturityStep2 = types.Int64Value(int64(result.Maturity.Step2))
	ps.MaturityStep3 = types.Int64Value(int64(result.Maturity.Step3))
	ps.MaturityStep4 = types.Int64Value(int64(result.Maturity.Step4))
	ps.MaturityStep5 = types.Int64Value(int64(result.Maturity.Step5))
	ps.Measures = getMeasuresFromMap(result.Measures)

	//Set refreshed state
	diags = resp.State.Set(ctx, &ps)
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

	//Get slices from plan
	dt := getSliceFromSetOfString(plan.DataTags)
	ct := getSliceFromSetOfString(plan.ComplianceTags)
	st := getSliceFromSetOfString(plan.SOCTags)
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

	measureMap := make(map[string]zerotrust.MeasureState, 0)

	//Loop through measures
	for k, m := range plan.Measures {

		//Check if measure exists
		if !sliceContains(r.getAvailableMeasures(), k) {
			resp.Diagnostics.AddError("Measure does not exists.",
				"Messure ["+k+"] does not exist, available measures ["+strings.Join(r.getAvailableMeasures(), ",")+"]")
			return
		}

		var assigned_timestamp int
		if !(m.Assigned_timestamp.IsUnknown() || m.Assigned_timestamp.IsNull()) {
			assigned_timestamp = int(m.Assigned_timestamp.ValueInt64())
		} else {
			assigned_timestamp = int(time.Now().Unix())
		}

		assignment := zerotrust.Assignment{
			Assigned:                 m.Assigned.ValueBool(),
			LastDeterminedByPersonID: m.Assigned_by.ValueString(),
			LastDeterminedTimestamp:  assigned_timestamp,
		}

		var implemented_timestamp int
		if !(m.Implemented_timestamp.IsUnknown() || m.Implemented_timestamp.IsNull()) {
			implemented_timestamp = int(m.Implemented_timestamp.ValueInt64())
		} else {
			implemented_timestamp = int(time.Now().Unix())
		}

		implementation := zerotrust.Implementation{
			Implemented:              m.Implemented.ValueBool(),
			LastDeterminedByPersonID: m.Implemented_by.ValueString(),
			LastDeterminedTimestamp:  implemented_timestamp,
		}

		var evidenced_timestamp int
		if !(m.Evidenced_timestamp.IsUnknown() || m.Evidenced_timestamp.IsNull()) {
			evidenced_timestamp = int(m.Evidenced_timestamp.ValueInt64())
		} else {
			evidenced_timestamp = int(time.Now().Unix())
		}

		evidence := zerotrust.Evidence{
			Evidenced:                m.Evidenced.ValueBool(),
			LastDeterminedByPersonID: m.Evidenced_by.ValueString(),
			LastDeterminedTimestamp:  evidenced_timestamp,
		}

		measureMap[k] = zerotrust.MeasureState{
			Assignment:     &assignment,
			Implementation: &implementation,
			Evidence:       &evidence,
		}
	}

	if len(measureMap) == 0 {
		measureMap = nil
	}
	protectsurface.Measures = measureMap

	result, err := r.client.ZeroTrust.UpdateProtectSurface(protectsurface)

	if err != nil {
		resp.Diagnostics.AddError("Error updating protect surface", "unexpected error: "+err.Error())
		return
	}

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
	plan.CustomerLabels, _ = types.MapValueFrom(ctx, types.StringType, result.CustomerLabels)
	plan.SOCTags = getSetOfStringFromSlice(result.SocTags)
	plan.AllowFlowsFromOutside = types.BoolValue(result.FlowsFromOutside.Allow)
	plan.AllowFlowsToOutside = types.BoolValue(result.FlowsToOutside.Allow)
	plan.MaturityStep1 = types.Int64Value(int64(result.Maturity.Step1))
	plan.MaturityStep2 = types.Int64Value(int64(result.Maturity.Step2))
	plan.MaturityStep3 = types.Int64Value(int64(result.Maturity.Step3))
	plan.MaturityStep4 = types.Int64Value(int64(result.Maturity.Step4))
	plan.MaturityStep5 = types.Int64Value(int64(result.Maturity.Step5))
	plan.Measures = getMeasuresFromMap(result.Measures)

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

func (r *protectsurfaceResource) getAvailableMeasures() []string {
	availableMeasures, _ := r.client.ZeroTrust.GetMeasures()
	availableMeasuresInSlice := make([]string, 0)
	for _, mg := range availableMeasures.Groups {
		for _, m := range mg.Measures {
			availableMeasuresInSlice = append(availableMeasuresInSlice, m.Name)
		}
	}

	return availableMeasuresInSlice
}

func getMeasuresFromMap(measureMap map[string]zerotrust.MeasureState) map[string]measure {
	if len(measureMap) == 0 {
		return nil
	}

	measures := make(map[string]measure, len(measureMap))

	for k, state := range measureMap {
		measures[k] = measure{
			//Measure:               types.StringValue(m),
			Assigned:              types.BoolValue(state.Assignment.Assigned),
			Assigned_by:           types.StringValue(state.Assignment.LastDeterminedByPersonID),
			Assigned_timestamp:    types.Int64Value(int64(state.Assignment.LastDeterminedTimestamp)),
			Implemented:           types.BoolValue(state.Implementation.Implemented),
			Implemented_by:        types.StringValue(state.Implementation.LastDeterminedByPersonID),
			Implemented_timestamp: types.Int64Value(int64(state.Implementation.LastDeterminedTimestamp)),
			Evidenced:             types.BoolValue(state.Evidence.Evidenced),
			Evidenced_by:          types.StringValue(state.Evidence.LastDeterminedByPersonID),
			Evidenced_timestamp:   types.Int64Value(int64(state.Evidence.LastDeterminedTimestamp)),
		}
	}

	return measures
}

// func getMapfromMeasures(measures []measure) map[string]zerotrust.MeasureState {
// 	measureMap := make(map[string]zerotrust.MeasureState, len(measures))

// 	for _, m := range measures {
// 		assignment := zerotrust.Assignment{
// 			Assigned:                 m.Assigned.ValueBool(),
// 			LastDeterminedByPersonID: m.Assigned_by.ValueString(),
// 			LastDeterminedTimestamp:  int(m.Assigned_timestamp.ValueInt64()),
// 		}

// 		implementation := zerotrust.Implementation{
// 			Implemented:              m.Implemented.ValueBool(),
// 			LastDeterminedByPersonID: m.Implemented_by.ValueString(),
// 			LastDeterminedTimestamp:  int(m.Implemented_timestamp.ValueInt64()),
// 		}

// 		evidence := zerotrust.Evidence{
// 			Evidenced:                m.Evidenced.ValueBool(),
// 			LastDeterminedByPersonID: m.Evidenced_by.ValueString(),
// 			LastDeterminedTimestamp:  int(m.Evidenced_timestamp.ValueInt64()),
// 		}

// 		measureMap[k] = zerotrust.MeasureState{
// 			Assignment:     &assignment,
// 			Implementation: &implementation,
// 			Evidence:       &evidence,
// 		}
// 	}

// 	return measureMap
// }
