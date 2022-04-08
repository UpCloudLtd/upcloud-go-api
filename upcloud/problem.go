package upcloud

import (
	"fmt"
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
	// CorrelationID is an unique string that identifies the request that caused the problem
	CorrelationID string `json:"correlation_id,omitempty"`
	// HTTP Status code
	Status int `json:"status"`
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

// ProblemInvalidParam is a type describing extra information in the Problem type's InvalidParams field.
type ProblemInvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}
