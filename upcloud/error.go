package upcloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/client"
)

func NewError(clientErr *client.Error) error {
	if clientErr == nil {
		return nil
	}

	ucErr := &Error{}

	err := json.Unmarshal(clientErr.ResponseBody, ucErr)
	if err != nil {
		return fmt.Errorf("received malformed client error: %s", err)
	}

	ucErr.Status = clientErr.ErrorCode
	ucErr.clientError = clientErr

	return ucErr
}

// Error represents an error
type Error struct {
	// Code is a short, programmatic identifier of an error
	Code string `json:"code"`

	// Message is a human-readable description of an error
	Message string `json:"message"`

	// Status is an HTTP Status code from the response
	Status int `json:"-"`

	clientError  *client.Error
	problemError *ProblemError
}

// IsProblem checks if the `Error` has additional information on the problem with HTTP request (conforming to RFC7807)
func (err *Error) IsProblem() bool {
	return err.clientError.Type == client.ErrorTypeProblem
}

// Problem returns additional information on the problem with HTTP request (conforming to RFC7807), assuming they were included in error response.
// Second return value is a boolean that will be true if the additional information actually exist.
func (err *Error) Problem() (*ProblemError, bool) {
	if err.IsProblem() {
		return err.problemError, true
	}

	return nil, false
}

// Error implements the Error interface
func (err *Error) Error() string {
	return fmt.Sprintf("%s (%s)", err.Message, err.Code)
}

func (err *Error) UnmarshalJSON(buf []byte) error {
	localError := struct {
		LegacyError *struct {
			ErrorCode    string `json:"error_code,omitempty"`
			ErrorMessage string `json:"error_message,omitempty"`
		} `json:"error,omitempty"`
		ProblemType          string                `json:"type,omitempty"`
		ProblemTitle         string                `json:"title,omitempty"`
		ProblemCorrelationID string                `json:"correlation_id,omitempty"`
		ProblemStatus        int                   `json:"status,omitempty"`
		ProblemInvalidParams []ProblemInvalidParam `json:"invalid_params,omitempty"`
	}{}

	unmarshalErr := json.Unmarshal(buf, &localError)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	// We got the legacy error struct from the API, just populate the main error fields
	if localError.LegacyError != nil {
		err.Code = localError.LegacyError.ErrorCode
		err.Message = localError.LegacyError.ErrorMessage
		return nil
	}

	// This means we got a json+problem error from the API, populate the fields accordingly
	err.problemError = &ProblemError{
		Status:        localError.ProblemStatus,
		Title:         localError.ProblemTitle,
		Type:          localError.ProblemType,
		InvalidParams: localError.ProblemInvalidParams,
		CorrelationID: localError.ProblemCorrelationID,
	}

	// We also need to populate main error fields
	err.Code = err.problemError.getShortenedType()
	err.Message = err.problemError.Title

	return nil
}

//// UnmarshalJSON is a custom unmarshaller that deals with
//// deeply embedded values.
//func (e *Error) UnmarshalJSON(b []byte) error {
//	type localError Error
//	v := struct {
//		Error localError `json:"error"`
//	}{}
//	err := json.Unmarshal(b, &v)
//	if err != nil {
//		return err
//	}
//
//	(*e) = Error(v.Error)
//
//	return nil
//}

// ProblemError is the type conforming to RFC7807 that represents an error or a problem associated with an HTTP request.
type ProblemError struct {
	// Type is the URI to a page describing the problem
	Type string `json:"type"`
	// Title is the human-readable description if the problem
	Title string `json:"title"`
	// InvalidParams if set, is a list of ProblemInvalidParam describing a specific part(s) of the request
	// that caused the problem
	InvalidParams []ProblemInvalidParam `json:"invalid_params,omitempty"`
	// CorrelationID is an unique string that identifies the request that caused the problem
	CorrelationID string `json:"correlation_id,omitempty"`
	// HTTP Status code
	Status int `json:"status"`
}

// Error implements the Error interface
func (pe *ProblemError) Error() string {
	var sb strings.Builder
	_, _ = fmt.Fprintf(&sb, "error: message=%q, type=%q", pe.Title, pe.Type)
	if pe.CorrelationID != "" {
		_, _ = fmt.Fprintf(&sb, ", correlation_id=%s", pe.CorrelationID)
	}
	if len(pe.InvalidParams) > 0 {
		for _, ip := range pe.InvalidParams {
			_, _ = fmt.Fprintf(&sb, ", invalid_params_%s='%s'", ip.Name, ip.Reason)
		}
	}
	return sb.String()
}

func (pe *ProblemError) getShortenedType() string {
	parts := strings.SplitN(pe.Type, "#", 2)

	if len(parts) < 2 {
		return ""
	}

	return strings.Replace(parts[1], "ERROR_", "", 1)
}

// ProblemErrorInvalidParam is a type describing extra information in the Problem type's InvalidParams field.
type ProblemErrorInvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}
