package auxo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// getSliceFromSetOfString converts a slice of basetypes.StringValue to a slice of string
func getSliceFromSetOfString(values []basetypes.StringValue) []string {
	var result []string
	for _, value := range values {
		result = append(result, value.String())
	}
	return result
}

func getSetOfStringFromSlice(values []string) []basetypes.StringValue {
	var result []basetypes.StringValue
	for _, value := range values {
		result = append(result, basetypes.NewStringValue(value))
	}
	return result
}

// TODO - This is part of the upcoming 1.2.0 release of the terraform-plugin-framework
// DefaultAttributePlanModifier is to set default value for an attribute
// https://github.com/hashicorp/terraform-plugin-framework/issues/285

func Int64DefaultValue(v types.Int64) planmodifier.Int64 {
	return &int64DefaultValuePlanModifier{v}
}

type int64DefaultValuePlanModifier struct {
	DefaultValue types.Int64
}

var _ planmodifier.Int64 = (*int64DefaultValuePlanModifier)(nil)

func (apm *int64DefaultValuePlanModifier) Description(ctx context.Context) string {
	/* ... */
	return "returns the default value"
}

func (apm *int64DefaultValuePlanModifier) MarkdownDescription(ctx context.Context) string {
	/* ... */
	return "returns the default value"
}

func (apm *int64DefaultValuePlanModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, res *planmodifier.Int64Response) {
	// If the attribute configuration is not null, we are done here
	if !req.ConfigValue.IsNull() {
		return
	}

	// If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
	// has already been applied, and we don't want to interfere.
	if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
		return
	}

	res.PlanValue = apm.DefaultValue
}

// --- End of DefaultAttributePlanModifier
