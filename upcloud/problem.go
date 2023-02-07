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

// MatchesProblemType is a helper method for comparing a `Problem.Type` field against a problem type constants defined in this package.
// Please note that while you can compare the `Type` field directly to those constants, it is not recommended and might not always work reliably.
// Examples:
//
//	res, err := svc.GetServerDetails(ctx, &req)
//
//	var problem *upcloud.Problem
//	if errors.As(err, &problem) {
//		// This is the recommended way
//		if problem.MatchesProblemType(upcloud.ProblemTypeResourceAlreadyExists) {
//			handleStuff(problem)
//		}
//
//		// This might work, but is not guaranteed and not recommended
//		if problem.Type == upcloud.ProblemTypeResourceAlreadyExists {
//			handleStuff(problem)
//		}
//	}
func (p *Problem) MatchesProblemType(problemType string) bool {
	// Type is a URL, we need to extract meaningful fragment from it for comparison purposes
	if strings.Contains(p.Type, "https://") {
		parts := strings.SplitN(p.Type, "#", 2)

		if len(parts) < 2 {
			return false
		}

		return strings.Replace(parts[1], "ERROR_", "", 1) == problemType
	}

	// Type is just bare problem type string (this happens when error returned from the API is not in the json+problem format)
	// We can just compare
	return p.Type == problemType
}
