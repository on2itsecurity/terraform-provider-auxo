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

func getAPIError(err error) apiError {
	var apiErr apiError
	//Workaround for API error messages from go-auxo
	cleanError := strings.Replace(err.Error(), "Not 200 or 201 ok, but 404, with body ", "", -1)
	json.Unmarshal([]byte(cleanError), &apiErr)
	return apiErr
}
