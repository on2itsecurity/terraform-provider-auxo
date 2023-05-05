package auxo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/on2itsecurity/go-auxo"
	"github.com/on2itsecurity/go-auxo/zerotrust"
)

var _ resource.Resource = &stateResource{}

type stateResource struct {
	client *auxo.Client
}

type stateResourceModel struct {
	ID             types.String   `tfsdk:"id"`
	Uniqueness_key types.String   `tfsdk:"uniqueness_key"`
	Description    types.String   `tfsdk:"description"`
	Protectsurface types.String   `tfsdk:"protectsurface_id"`
	Location       types.String   `tfsdk:"location_id"`
	ContentType    types.String   `tfsdk:"content_type"`
	ExistsOnAssets []types.String `tfsdk:"exists_on_assets"`
	Maintainer     types.String   `tfsdk:"maintainer"`
	Content        []types.String `tfsdk:"content"`
}

func NewStateResource() resource.Resource {
	return &stateResource{}
}

func (r *stateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_state"
}

func (r *stateResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	// Retrieve the client from the provider config
	r.client = req.ProviderData.(*auxo.Client)
}

func (r *stateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Auxo State",
		MarkdownDescription: "Auxo State",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Computed unique ID of the resource state",
				MarkdownDescription: "Computed unique ID of the resource state",
				Required:            false,
				Optional:            false,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"uniqueness_key": schema.StringAttribute{
				Description:         "Custom and optinal uniqueness key to identify the resource state",
				MarkdownDescription: "Custom and optinal uniqueness key to identify the resource state",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description:         "Description of the resource state",
				MarkdownDescription: "Description of the resource state",
				Required:            true,
			},
			"protectsurface_id": schema.StringAttribute{
				Description:         "ID of the protect surface",
				MarkdownDescription: "ID of the protect surface",
				Required:            true,
			},
			"location_id": schema.StringAttribute{
				Description:         "ID of the location",
				MarkdownDescription: "ID of the location",
				Required:            true,
			},
			"content_type": schema.StringAttribute{
				Description:         "Content type of the state i.e. ipv4, ipv6, azure_resource",
				MarkdownDescription: "Content type of the state i.e. ipv4, ipv6, azure_resource",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("ipv4"),
			},
			"exists_on_assets": schema.SetAttribute{
				Description:         "Contains asset IDs which could match this state",
				MarkdownDescription: "Contains asset IDs which could match this state",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"maintainer": schema.StringAttribute{
				Description:         "Maintainer of the state either api or portal_manual",
				MarkdownDescription: "Maintainer of the state either api or portal_manual",
				Optional:            true,
			},
			"content": schema.SetAttribute{
				Description:         "Content of the state e.g. \"10.1.1.2/32\",\"10.1.1.3/32\"",
				MarkdownDescription: "Content of the state e.g. \"10.1.1.2/32\",\"10.1.1.3/32\"",
				Required:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *stateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan stateResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get content
	var content []string
	for _, c := range plan.Content {
		content = append(content, c.String())
	}

	// Get assets
	var assets []string
	for _, a := range plan.ExistsOnAssets {
		assets = append(assets, a.String())
	}

	// Create state (object)
	state := zerotrust.State{
		ID:               plan.ID.ValueString(),
		UniquenessKey:    plan.Uniqueness_key.ValueString(),
		Description:      plan.Description.ValueString(),
		ProtectSurface:   plan.Protectsurface.ValueString(),
		Location:         plan.Location.ValueString(),
		ContentType:      plan.ContentType.ValueString(),
		ExistsOnAssetIDs: assets,
		Maintainer:       plan.Maintainer.ValueString(),
		Content:          &content,
	}

	// Create state (API)
	result, err := r.client.ZeroTrust.CreateStateByObject(state)

	if err != nil {
		resp.Diagnostics.AddError("Error creating state", "unexpected error: "+err.Error())
	}

	// Map resonse to schema
	plan.ID = types.StringValue(result.ID)
	plan.Uniqueness_key = types.StringValue(result.UniquenessKey)
	plan.Description = types.StringValue(result.Description)
	plan.Protectsurface = types.StringValue(result.ProtectSurface)
	plan.Location = types.StringValue(result.Location)
	plan.ContentType = types.StringValue(result.ContentType)
	for _, a := range result.ExistsOnAssetIDs {
		plan.ExistsOnAssets = append(plan.ExistsOnAssets, types.StringValue(a))
	}
	plan.Maintainer = types.StringValue(result.Maintainer)
	for _, c := range *result.Content {
		plan.Content = append(plan.Content, types.StringValue(c))
	}

	// Set state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

}
func (r *stateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state stateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed state from AUXO
	result, err := r.client.ZeroTrust.GetStateByID(state.ID.ValueString())
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

	//Overwrite state with refreshed state
	//ID and UK cannot have changed
	state.Description = types.StringValue(result.Description)
	state.Protectsurface = types.StringValue(result.ProtectSurface)
	state.Location = types.StringValue(result.Location)
	state.ContentType = types.StringValue(result.ContentType)
	for _, a := range result.ExistsOnAssetIDs {
		state.ExistsOnAssets = append(state.ExistsOnAssets, types.StringValue(a))
	}
	state.Maintainer = types.StringValue(result.Maintainer)
	for _, c := range *result.Content {
		state.Content = append(state.Content, types.StringValue(c))
	}

	//Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
func (r *stateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Retrieve values from plan
	var plan stateResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get content
	var content []string
	for _, c := range plan.Content {
		content = append(content, c.String())
	}

	// Get assets
	var assets []string
	for _, a := range plan.ExistsOnAssets {
		assets = append(assets, a.String())
	}

	// Create state (object)
	state := zerotrust.State{
		ID:               plan.ID.ValueString(),
		UniquenessKey:    plan.Uniqueness_key.ValueString(),
		Description:      plan.Description.ValueString(),
		ProtectSurface:   plan.Protectsurface.ValueString(),
		Location:         plan.Location.ValueString(),
		ContentType:      plan.ContentType.ValueString(),
		ExistsOnAssetIDs: assets,
		Maintainer:       plan.Maintainer.ValueString(),
		Content:          &content,
	}

	// Create state (API)
	result, err := r.client.ZeroTrust.CreateStateByObject(state)

	if err != nil {
		resp.Diagnostics.AddError("Error creating state", "unexpected error: "+err.Error())
	}

	// Map resonse to schema
	plan.ID = types.StringValue(result.ID)
	plan.Uniqueness_key = types.StringValue(result.UniquenessKey)
	plan.Description = types.StringValue(result.Description)
	plan.Protectsurface = types.StringValue(result.ProtectSurface)
	plan.Location = types.StringValue(result.Location)
	plan.ContentType = types.StringValue(result.ContentType)
	for _, a := range result.ExistsOnAssetIDs {
		plan.ExistsOnAssets = append(plan.ExistsOnAssets, types.StringValue(a))
	}
	plan.Maintainer = types.StringValue(result.Maintainer)
	for _, c := range *result.Content {
		plan.Content = append(plan.Content, types.StringValue(c))
	}

	// Set state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

}
func (r *stateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state stateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete state
	err := r.client.ZeroTrust.DeleteStateByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting state", "unexpected error: "+err.Error())
		return
	}
}
