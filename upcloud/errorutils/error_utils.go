package errorutils

import (
	"net/http"
	"strings"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
)

// IsNotFoundError checks if provided error is API Not Found error
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	if uError, ok := err.(*upcloud.Error); ok {
		return uError.Status == http.StatusNotFound
	}

	if uProblem, ok := err.(*upcloud.Problem); ok {
		return uProblem.Status == http.StatusNotFound
	}

	return false
}

// List of all error codes that indicate that API object already exists
var alreadyExistsErrCodes = map[string]struct{}{
	ErrCodeFirewallRuleExists: {},
	ErrCodeDBExists:           {},
	ErrCodeServiceExists:      {},
	ErrCodeInterfaceExists:    {},
	ErrCodeTagExists:          {},
}

// List of all json+problem type fragments that indicate that API object already exists
var alreadyExistsJsonProblemTypes = map[string]struct{}{
	ProblemTypeResourceAlreadyExists: {},
}

// IsAlreadyExistsError checks if provided error is API Already Exists error
func IsAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}

	if uError, ok := err.(*upcloud.Error); ok {
		_, errCodeOk := alreadyExistsErrCodes[uError.ErrorCode]
		return uError.Status == http.StatusConflict && errCodeOk
	}

	if uProblem, ok := err.(*upcloud.Problem); ok {
		problemType := getJsonProblemType(uProblem)
		_, problemTypeOk := alreadyExistsJsonProblemTypes[problemType]
		return uProblem.Status == http.StatusBadRequest && problemTypeOk
	}

	return false
}

func IsAuthenticationFailedError(err error) bool {
	if err == nil {
		return false
	}

	if uError, ok := err.(*upcloud.Error); ok {
		return uError.Status == http.StatusUnauthorized
	}

	if uProblem, ok := err.(*upcloud.Problem); ok {
		return uProblem.Status == http.StatusUnauthorized
	}

	return false
}

// GetJsonProblemType gets the meaningful part of json+problem type field
// json+problem `type` field should be a URL to a page that explains the error
// for the lack of better alternatives we need to use a fragment of that URL for comparing purposes
func getJsonProblemType(err *upcloud.Problem) string {
	parts := strings.SplitN(err.Type, "#", 2)

	if len(parts) < 2 {
		return ""
	}

	return parts[1]
}
