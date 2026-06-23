package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// LogFn is a function that logs a message with context and optional key-value pairs, e.g., slog.DebugContext
type LogFn func(context.Context, string, ...any)

type Logger struct {
	logFn LogFn
}

func NewLogger(logFn LogFn) *Logger {
	return &Logger{logFn: logFn}
}

func (l *Logger) LogRequest(r *http.Request, body []byte) {
	const authorization string = "Authorization"

	if l != nil && l.logFn != nil {
		headers := r.Header.Clone()
		if _, ok := headers[authorization]; ok {
			auth := strings.Split(headers.Get(authorization), " ")
			// Redact the token part of the Authorization header or the whole value if there is no space to separate scheme from parameters.
			if len(auth) > 1 {
				headers.Set(authorization, fmt.Sprintf("%s [REDACTED]", auth[0]))
			} else {
				headers.Set(authorization, "[REDACTED]")
			}
		}

		l.logFn(r.Context(), "Sending request to UpCloud API",
			"url", r.URL.Redacted(),
			"method", r.Method,
			"headers", headers,
			"body", prettyJSON(body),
		)
	}
}

func (l *Logger) LogResponse(r *http.Response, body []byte) {
	if l != nil && l.logFn != nil {
		l.logFn(r.Request.Context(), "Received response from UpCloud API",
			"url", r.Request.URL.Redacted(),
			"status", r.Status,
			"headers", r.Header,
			"body", prettyJSON(body),
		)
	}
}

// Pretty prints given JSON bytes. If the JSON is not valid, returns the original bytes as string.
func prettyJSON(i []byte) string {
	var o bytes.Buffer
	if err := json.Indent(&o, i, "", "  "); err != nil {
		return string(i)
	}
	return o.String()
}
