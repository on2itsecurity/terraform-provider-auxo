package auxo

import (
	"context"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/on2itsecurity/go-auxo/v2"
	"github.com/on2itsecurity/go-auxo/v2/zerotrust"
)

var _ resource.Resource = &stateResource{}

type stateResource struct {
	client *auxo.Client
	mutex  *sync.Mutex
}

type stateResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Uniqueness_key types.String `tfsdk:"uniqueness_key"`
	Description    types.String `tfsdk:"description"`
	Protectsurface types.String `tfsdk:"protectsurface_id"`
	Location       types.String `tfsdk:"location_id"`
	ContentType    types.String `tfsdk:"content_type"`
	ExistsOnAssets types.Set    `tfsdk:"exists_on_assets"`
	Maintainer     types.String `tfsdk:"maintainer"`
	Content        types.Set    `tfsdk:"content"`
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
	c := req.ProviderData.(*auxoClient)
	r.client = c.client
	r.mutex = c.m
}

func (r *stateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "A state contains resources and their location, belonging to a protect surface.",
		MarkdownDescription: "A state contains resources and their location, belonging to a protect surface.",
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
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"maintainer": schema.StringAttribute{
				Description:         "Maintainer of the state either api or portal_manual",
				MarkdownDescription: "Maintainer of the state either api or portal_manual",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("api_terraform"),
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

	// Create state (object)
	state := resourceModelToState(&plan, ctx)

	// Create state (API)
	result, err := r.client.ZeroTrust.CreateStateByObject(ctx, state)

	if err != nil {
		resp.Diagnostics.AddError("Error creating state", "unexpected error: "+err.Error())
		return
	}

	// Map resonse to schema
	plan = stateToResourceModel(result, ctx)

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
	result, err := r.client.ZeroTrust.GetStateByID(ctx, state.ID.ValueString())
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
	state = stateToResourceModel(result, ctx)

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

	// Create state (object)
	state := resourceModelToState(&plan, ctx)

	// Create state (API)
	result, err := r.client.ZeroTrust.CreateStateByObject(ctx, state)

	if err != nil {
		resp.Diagnostics.AddError("Error creating state", "unexpected error: "+err.Error())
	}

	// Map resonse to schema
	plan = stateToResourceModel(result, ctx)

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
	err := r.client.ZeroTrust.DeleteStateByID(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting state", "unexpected error: "+err.Error())
		return
	}
}

// resourceModelToState maps the resource model to the zerotrust.state object
func resourceModelToState(m *stateResourceModel, ctx context.Context) zerotrust.State {
	var existsOnAssets, content []string

	if !m.ExistsOnAssets.IsNull() {
		_ = m.ExistsOnAssets.ElementsAs(ctx, &existsOnAssets, false)
	}

	if !m.Content.IsNull() {
		_ = m.Content.ElementsAs(ctx, &content, false)
	}

	state := zerotrust.State{
		ID:               m.ID.ValueString(),
		UniquenessKey:    m.Uniqueness_key.ValueString(),
		Description:      m.Description.ValueString(),
		ProtectSurface:   m.Protectsurface.ValueString(),
		Location:         m.Location.ValueString(),
		ContentType:      m.ContentType.ValueString(),
		ExistsOnAssetIDs: existsOnAssets,
		Maintainer:       m.Maintainer.ValueString(),
		Content:          &content,
	}
	return state
}

// StateToResouceModel maps the zerotrust.state object to the resource model
func stateToResourceModel(state *zerotrust.State, ctx context.Context) stateResourceModel {
	var existsOnAssets, content basetypes.SetValue

	if state.ExistsOnAssetIDs != nil {
		existsOnAssets, _ = types.SetValueFrom(ctx, types.StringType, state.ExistsOnAssetIDs)
	}
	if state.Content != nil {
		content, _ = types.SetValueFrom(ctx, types.StringType, *state.Content)
	}

	return stateResourceModel{
		ID:             types.StringValue(state.ID),
		Uniqueness_key: types.StringValue(state.UniquenessKey),
		Description:    types.StringValue(state.Description),
		Protectsurface: types.StringValue(state.ProtectSurface),
		Location:       types.StringValue(state.Location),
		ContentType:    types.StringValue(state.ContentType),
		ExistsOnAssets: existsOnAssets,
		Maintainer:     types.StringValue(state.Maintainer),
		Content:        content,
	}
}
