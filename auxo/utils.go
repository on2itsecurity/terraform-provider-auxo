package auxo

import (
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// apiError is the struct for the error returned by the go-auxo API
type apiError struct {
	ID      string `json:"error_id"`
	Name    string `json:"error_name"`
	Message string `json:"error_message"`
}

// getAPIError returns an apiError struct from a go-auxo error
func getAPIError(err error) apiError {
	var apiErr apiError
	//Workaround for API error messages from go-auxo
	cleanError := strings.Replace(err.Error(), "Not 200 or 201 ok, but 404, with body ", "", -1)
	json.Unmarshal([]byte(cleanError), &apiErr)
	return apiErr
}

// getSliceFromSetOfString converts a slice of basetypes.StringValue to a slice of string
func getSliceFromSetOfString(values []basetypes.StringValue) []string {
	var result []string
	for _, value := range values {
		result = append(result, value.ValueString())
	}
	return result
}

// getSetOfStringFromSlice converts a slice of string to a slice of basetypes.StringValue
func getSetOfStringFromSlice(values []string) []basetypes.StringValue {
	var result []basetypes.StringValue
	for _, value := range values {
		result = append(result, basetypes.NewStringValue(value))
	}
	return result
}

// func convertTFMapToGo(m map[types.String]types.String) map[string]string {
// 	if len(m) == 0 {
// 		return nil
// 	}
// 	result := make(map[string]string)
// 	for k, v := range m {
// 		result[k.ValueString()] = v.ValueString()
// 	}
// 	return result
// }

// func convertGoMapToTF(m map[string]string) map[types.String]types.String {
// 	if len(m) == 0 {
// 		return nil
// 	}
// 	result := make(map[types.String]types.String)
// 	for k, v := range m {
// 		result[basetypes.NewStringValue(k)] = basetypes.NewStringValue(v)
// 	}
// 	return result
// }

// sliceContains checks if a string is in a slice of strings
func sliceContains(slice []string, match string) bool {
	for _, str := range slice {
		if str == match {
			return true
		}
	}
	return false
}
