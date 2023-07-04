package auxo

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/on2itsecurity/go-auxo"
	"github.com/on2itsecurity/go-auxo/zerotrust"
)

var _ resource.Resource = &measureResource{}

type measureResource struct {
	client *auxo.Client
}

type measureResourceModel struct {
	Protectsurface types.String       `tfsdk:"protectsurface"`
	Measures       map[string]measure `tfsdk:"measures"`
}

type measure struct {
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

func NewMeasureResource() resource.Resource {
	return &measureResource{}
}

func (r *measureResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_measure"
}

func (r *measureResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	// Retrieve the client from the provider config
	r.client = req.ProviderData.(*auxo.Client)
}

func (r *measureResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "A Measure resource represents the measures for the specified protectsurface.",
		MarkdownDescription: "A Measure resource represents the measures for the specified protectsurface.",
		Attributes: map[string]schema.Attribute{
			"protectsurface": schema.StringAttribute{
				Description:         "The ID of the protectsurface",
				MarkdownDescription: "The ID of the protectsurface",
				Required:            true,
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

func (r *measureResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan measureResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get PS and add measures
	ps, diags := r.resourceModelToCompletePS(&plan, ctx)

	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Create(=update) PS
	ps, err := r.client.ZeroTrust.UpdateProtectSurface(*ps)

	if err != nil {
		resp.Diagnostics.AddError("Error creating protectsurface", "unexpected error: "+err.Error())
		return
	}

	//Read back the measures
	measures := getMeasuresFromMap(ps.Measures)
	plan.Protectsurface = types.StringValue(ps.ID)
	plan.Measures = measures

	// Set state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *measureResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state measureResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed state from AUXO
	result, err := r.client.ZeroTrust.GetProtectSurfaceByID(state.Protectsurface.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading measures", "unexpected error: "+err.Error())
		return
	}

	//Overwrite state with refreshed state
	measures := getMeasuresFromMap(result.Measures)
	state.Protectsurface = types.StringValue(result.ID)
	state.Measures = measures

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *measureResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan measureResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get PS and add measures
	ps, diags := r.resourceModelToCompletePS(&plan, ctx)

	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Create(=update) PS
	ps, err := r.client.ZeroTrust.UpdateProtectSurface(*ps)

	if err != nil {
		resp.Diagnostics.AddError("Error creating protectsurface", "unexpected error: "+err.Error())
		return
	}

	//Read back the measures
	measures := getMeasuresFromMap(ps.Measures)
	plan.Protectsurface = types.StringValue(ps.ID)
	plan.Measures = measures

	// Set state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *measureResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state measureResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get PS and remove flows
	ps, err := r.client.ZeroTrust.GetProtectSurfaceByID(state.Protectsurface.ValueString())
	ps.Measures = map[string]zerotrust.MeasureState{}

	// Update PS, with deleted flows
	_, err = r.client.ZeroTrust.CreateProtectSurfaceByObject(*ps, true)

	if err != nil {
		resp.Diagnostics.AddError("Error deleting measures", "unexpected error: "+err.Error())
		return
	}
}

func (r *measureResource) getAvailableMeasures() []string {
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

func (r *measureResource) resourceModelToCompletePS(plan *measureResourceModel, ctx context.Context) (*zerotrust.ProtectSurface, diag.Diagnostics) {
	var diags diag.Diagnostics

	psID := plan.Protectsurface.ValueString()
	ps, err := r.client.ZeroTrust.GetProtectSurfaceByID(psID)

	if err != nil {
		diags.AddError("Error getting protectsurface", "unexpected error: "+err.Error())
		return nil, diags
	}

	measureMap := make(map[string]zerotrust.MeasureState, 0)

	//Loop through measures
	for k, m := range plan.Measures {

		//Check if measure exists
		if !sliceContains(r.getAvailableMeasures(), k) {
			diags.AddError("Measure does not exists.",
				"Messure ["+k+"] does not exist, available measures ["+strings.Join(r.getAvailableMeasures(), ",")+"]")
			return nil, diags
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
	ps.Measures = measureMap

	return ps, diags
}
