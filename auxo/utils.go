package auxo

import (
	"encoding/json"
	"strings"
)

type apiError struct {
	ID      string `json:"error_id"`
	Name    string `json:"error_name"`
	Message string `json:"error_message"`
}

// createStringSliceFromListInput converts a slice/list of interface{} to a slice of strings
func createStringSliceFromListInput(inputList []interface{}) []string {
	output := make([]string, len(inputList))
	for k, v := range inputList {
		output[k] = v.(string)
	}

	return output
}

// getAPIError returns an apiError struct from a go-auxo error
func getAPIError(err error) apiError {
	var apiErr apiError
	//Workaround for API error messages from go-auxo
	cleanError := strings.Replace(err.Error(), "Not 200 or 201 ok, but 404, with body ", "", -1)
	json.Unmarshal([]byte(cleanError), &apiErr)
	return apiErr
}

// sliceContains checks if a string is in a slice of strings
func sliceContains(slice []string, match string) bool {
	for _, str := range slice {
		if str == match {
			return true
		}
	}
	return false
}
