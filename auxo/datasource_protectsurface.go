package auxo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/on2itsecurity/go-auxo"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &protectsurfaceDataSource{}
	_ datasource.DataSourceWithConfigure = &protectsurfaceDataSource{}
)

type protectsurfaceDataSource struct {
	client *auxo.Client
}

type protectsurfaceDataSourceModel struct {
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

// NewprotectsurfaceDataSource is a helper function to simplify the provider implementation.
func NewProtectsurfaceDataSource() datasource.DataSource {
	return &protectsurfaceDataSource{}
}

// Metadata returns the data source type name.
func (d *protectsurfaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_protectsurface"
}

func (d *protectsurfaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	// Retrieve the client from the provider config
	d.client = req.ProviderData.(*auxo.Client)
}

// Schema defines the schema for the data source.
func (d *protectsurfaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "A protectsurface which can be used in a state object to reflect the protectsurface of the resources. Either by specifying the name or the uniqueness_key",
		MarkdownDescription: "A protectsurface which can be used in a state object to reflect the protectsurface of the resources. Either by specifying the name or the uniqueness_key",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Computed unique IDs of the protectsurface",
				MarkdownDescription: "Computed unique IDs of the protectsurface",
				Computed:            true,
			},
			"uniqueness_key": schema.StringAttribute{
				Description:         "Uniqueness key of the protectsurface",
				MarkdownDescription: "Uniqueness key of the protectsurface",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the protectsurface",
				MarkdownDescription: "Name of the protectsurface",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the resource protectsurface",
				MarkdownDescription: "Description of the resource protectsurface",
				Computed:            true,
			},
			"main_contact": schema.StringAttribute{
				Description:         "Main contact of the resource protectsurface",
				MarkdownDescription: "Main contact of the resource protectsurface",
				Computed:            true,
			},
			"security_contact": schema.StringAttribute{
				Description:         "Security contact of the resource protectsurface",
				MarkdownDescription: "Security contact of the resource protectsurface",
				Computed:            true,
			},
			"in_control_boundary": schema.BoolAttribute{
				Description:         "This protect surface is within the 'control boundary'",
				MarkdownDescription: "This protect surface is within the 'control boundary'",
				Computed:            true,
			},
			"in_zero_trust_focus": schema.BoolAttribute{
				Description:         "This protect surface is within the 'zero trust focus' (actively maintained and monitored)",
				MarkdownDescription: "This protect surface is within the 'zero trust focus' (actively maintained and monitored)",
				Computed:            true,
			},
			"relevance": schema.Int64Attribute{
				Description:         "Relevance of the resource protectsurface",
				MarkdownDescription: "Relevance of the resource protectsurface",
				Computed:            true,
			},
			"confidentiality": schema.Int64Attribute{
				Description:         "Confidentiality of the resource protectsurface",
				MarkdownDescription: "Confidentiality of the resource protectsurface",
				Computed:            true,
			},
			"integrity": schema.Int64Attribute{
				Description:         "Integrity of the resource protectsurface",
				MarkdownDescription: "Integrity of the resource protectsurface",
				Computed:            true,
			},
			"availability": schema.Int64Attribute{
				Description:         "Availability of the resource protectsurface",
				MarkdownDescription: "Availability of the resource protectsurface",
				Computed:            true,
			},
			"data_tags": schema.SetAttribute{
				Description:         "Data tags of the resource protectsurface",
				MarkdownDescription: "Data tags of the resource protectsurface",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"compliance_tags": schema.SetAttribute{
				Description:         "Compliance tags of the resource protectsurface",
				MarkdownDescription: "Compliance tags of the resource protectsurface",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"customer_labels": schema.MapAttribute{
				Description:         "Customer labels of the resource protectsurface",
				MarkdownDescription: "Customer labels of the resource protectsurface",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"soc_tags": schema.SetAttribute{
				Description:         "SOC tags of the resource protectsurface, only use when advised by the SOC",
				MarkdownDescription: "SOC tags of the resource protectsurface, only use when advised by the SOC",
				Computed:            true,
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
				Computed:            true,
			},
			"maturity_step2": schema.Int64Attribute{
				Description:         "Maturity step 2",
				MarkdownDescription: "Maturity step 2",
				Computed:            true,
			},
			"maturity_step3": schema.Int64Attribute{
				Description:         "Maturity step 3",
				MarkdownDescription: "Maturity step 3",
				Computed:            true,
			},
			"maturity_step4": schema.Int64Attribute{
				Description:         "Maturity step 4",
				MarkdownDescription: "Maturity step 4",
				Computed:            true,
			},
			"maturity_step5": schema.Int64Attribute{
				Description:         "Maturity step 5",
				MarkdownDescription: "Maturity step 5",
				Computed:            true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *protectsurfaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state protectsurfaceDataSourceModel

	//Get input
	var input protectsurfaceDataSourceModel
	diags := req.Config.Get(ctx, &input)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Get protectsurfaces
	protectsurfaces, err := d.client.ZeroTrust.GetProtectSurfaces()
	if err != nil {
		resp.Diagnostics.AddError("Unable to retrieve protectsurfaces", err.Error())
		return
	}

	//Make map to check for duplicates (which allowed, but makes it impossible to look for a protectsurface by name)
	psCount := make(map[string]int)
	for _, protectsurface := range protectsurfaces {
		psCount[protectsurface.Name]++
	}
	if psCount[input.Name.ValueString()] > 1 {
		resp.Diagnostics.AddError("Duplicate name on backend, please use uniqueness key", "Duplicate name on backend name "+input.Name.ValueString()+", please use uniqueness key")
		return
	}

	//Check if one of the two is set
	if (input.Uniqueness_key.IsNull() && input.Name.IsNull()) || (!input.Uniqueness_key.IsNull() && !input.Name.IsNull()) {
		resp.Diagnostics.AddError("Either uniqueness_key OR name must be set", "")
		return
	}

	//Find the protectsurface
	for _, ps := range protectsurfaces {
		if (ps.UniquenessKey == input.Uniqueness_key.ValueString()) || (ps.Name == input.Name.ValueString()) {
			state.ID = types.StringValue(ps.ID)
			state.Name = types.StringValue(ps.Name)
			state.Uniqueness_key = types.StringValue(ps.UniquenessKey)
			//Map all "other" fields
			cl, _ := types.MapValueFrom(ctx, types.StringType, ps.CustomerLabels)

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

			state.Description = types.StringValue(ps.Description)
			state.MainContact = types.StringValue(ps.MainContactPersonID)
			state.SecurityContact = types.StringValue(ps.SecurityContactPersonID)
			state.InControlBoundary = types.BoolValue(ps.InControlBoundary)
			state.InZeroTrustFocus = types.BoolValue(ps.InZeroTrustFocus)
			state.Relevance = types.Int64Value(int64(ps.Relevance))
			state.Confidentiality = types.Int64Value(int64(ps.Confidentiality))
			state.Integrity = types.Int64Value(int64(ps.Integrity))
			state.Availability = types.Int64Value(int64(ps.Availability))
			state.DataTags = dt
			state.ComplianceTags = ct
			state.CustomerLabels = cl
			state.SOCTags = st
			state.AllowFlowsFromOutside = types.BoolPointerValue(ps.FlowsFromOutside.Allow)
			state.AllowFlowsToOutside = types.BoolPointerValue(ps.FlowsToOutside.Allow)
			state.MaturityStep1 = types.Int64Value(int64(ps.Maturity.Step1))
			state.MaturityStep2 = types.Int64Value(int64(ps.Maturity.Step2))
			state.MaturityStep3 = types.Int64Value(int64(ps.Maturity.Step3))
			state.MaturityStep4 = types.Int64Value(int64(ps.Maturity.Step4))
			state.MaturityStep5 = types.Int64Value(int64(ps.Maturity.Step5))

			break
		}
	}

	if state.ID.IsNull() {
		resp.Diagnostics.AddError("Unable to find protectsurface", "Unable to find protectsurface with uniqueness_key "+input.Uniqueness_key.ValueString()+" or name "+input.Name.ValueString())
		return
	}

	//set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
