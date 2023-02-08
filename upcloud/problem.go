package upcloud

import (
	"fmt"
	"net/url"
	"strings"
)

// Problem is the type conforming to RFC7807 that represents an error or a problem associated with an HTTP request.
type Problem struct {
	// Type is the URI to a page describing the problem
	Type string `json:"type"`
	// Title is the human-readable description if the problem
	Title string `json:"title"`
	// InvalidParams if set, is a list of ProblemInvalidParam describing a specific part(s) of the request
	// that caused the problem
	InvalidParams []ProblemInvalidParam `json:"invalid_params,omitempty"`
	// CorrelationID is a unique string that identifies the request that caused the problem
	// Please note that it is not always available
	CorrelationID string `json:"correlation_id,omitempty"`
	// HTTP Status code
	Status int `json:"status"`
}

// ProblemInvalidParam is a type describing extra information in the Problem type's InvalidParams field.
type ProblemInvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

func (p *Problem) Error() string {
	var sb strings.Builder
	_, _ = fmt.Fprintf(&sb, "error: message=%q, type=%q", p.Title, p.Type)
	if p.CorrelationID != "" {
		_, _ = fmt.Fprintf(&sb, ", correlation_id=%s", p.CorrelationID)
	}
	if len(p.InvalidParams) > 0 {
		for _, ip := range p.InvalidParams {
			_, _ = fmt.Fprintf(&sb, ", invalid_params_%s='%s'", ip.Name, ip.Reason)
		}
	}
	return sb.String()
}

// ErrorCode returns a short string that identifies the error; it should be used for programmatic comparisons
func (p *Problem) ErrorCode() string {
	// Type is a URL, we need to extract meaningful fragment from it for comparison purposes
	if strings.HasPrefix(p.Type, "https://") {
		parsedURL, err := url.Parse(p.Type)
		if err != nil {
			return ""
		}

		return strings.Replace(parsedURL.Fragment, "ERROR_", "", 1)
	}

	// Type is just bare error code string (this happens when error returned from the API is not in the json+problem format)
	return p.Type
}
