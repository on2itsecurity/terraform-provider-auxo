package auxo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/on2itsecurity/go-auxo"
	"github.com/on2itsecurity/go-auxo/zerotrust"
)

// Ensure the implementation satisfies the resource.Resource interface.
var _ resource.Resource = &locationResource{}

type locationResource struct {
	client *auxo.Client
}

type locationResourceModel struct {
	ID             types.String  `tfsdk:"id"`
	Uniqueness_key types.String  `tfsdk:"uniqueness_key"`
	Name           types.String  `tfsdk:"name"`
	Latitude       types.Float64 `tfsdk:"latitude"`
	Longitude      types.Float64 `tfsdk:"longitude"`
}

func NewLocationResource() resource.Resource {
	return &locationResource{}
}

func (r *locationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_location"
}

func (r *locationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	// Retrieve the client from the provider config
	r.client = req.ProviderData.(*auxo.Client)
}

func (r *locationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Auxo Location",
		MarkdownDescription: "Auxo Location",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Computed unique ID of the resource location",
				MarkdownDescription: "Computed unique ID of the resource location",
				Required:            false,
				Optional:            false,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"uniqueness_key": schema.StringAttribute{
				Description:         "Custom and optinal uniqueness key to identify the resource location",
				MarkdownDescription: "Custom and optinal uniqueness key to identify the resource location",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the resource location",
				MarkdownDescription: "Name of the resource location",
				Required:            true,
			},
			"latitude": schema.Float64Attribute{
				Description:         "Latitude of the resource location",
				MarkdownDescription: "Latitude of the resource location",
				Optional:            true,
				Computed:            true,
				Default:             float64default.StaticFloat64(0),
			},
			"longitude": schema.Float64Attribute{
				Description:         "Longitude of the resource location",
				MarkdownDescription: "Longitude of the resource location",
				Optional:            true,
				Computed:            true,
				Default:             float64default.StaticFloat64(0),
			},
		},
	}
}

func (r *locationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan locationResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create location (object)
	location := resourceModelToLocation(&plan)

	// Create location (API)
	result, err := r.client.ZeroTrust.CreateLocationByObject(location, false)

	if err != nil {
		resp.Diagnostics.AddError("Error creating location", "unexpected error: "+err.Error())
		return
	}

	// Map resonse to schema
	plan = locationToResourceModel(result)

	// Set state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *locationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var location locationResourceModel
	diags := req.State.Get(ctx, &location)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed location from AUXO
	result, err := r.client.ZeroTrust.GetLocationByID(location.ID.ValueString())
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

	//Overwrite state with refreshed location
	location = locationToResourceModel(result)

	//Set refreshed state
	diags = resp.State.Set(ctx, &location)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *locationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Retrieve values from plan
	var plan, state locationResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create location (object)
	location := resourceModelToLocation(&plan)

	// Update location (API)
	result, err := r.client.ZeroTrust.UpdateLocation(location)

	if err != nil {
		resp.Diagnostics.AddError("Error updating location", "unexpected error: "+err.Error())
		return
	}

	// Update state
	plan = locationToResourceModel(result)
	diags = resp.State.Set(ctx, plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *locationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//Retrieve values from state
	var location locationResourceModel
	diags := req.State.Get(ctx, &location)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Delete location
	err := r.client.ZeroTrust.DeleteLocationByID(location.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting location", "unexpected error: "+err.Error())
		return
	}
}

// resourceModelToLocation maps the resource model to the zerotrust.location object
func resourceModelToLocation(m *locationResourceModel) zerotrust.Location {
	return zerotrust.Location{
		ID:            m.ID.ValueString(),
		UniquenessKey: m.Uniqueness_key.ValueString(),
		Name:          m.Name.ValueString(),
		Coords: zerotrust.Coords{
			Latitude:  m.Latitude.ValueFloat64(),
			Longitude: m.Longitude.ValueFloat64(),
		},
	}
}

// locationToResouceModel maps the zerotrust.location object to the resource model
func locationToResourceModel(location *zerotrust.Location) locationResourceModel {
	return locationResourceModel{
		ID:             types.StringValue(location.ID),
		Uniqueness_key: types.StringValue(location.UniquenessKey),
		Name:           types.StringValue(location.Name),
		Latitude:       types.Float64Value(location.Coords.Latitude),
		Longitude:      types.Float64Value(location.Coords.Longitude),
	}
}
