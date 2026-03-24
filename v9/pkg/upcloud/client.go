package upcloud

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"os"

	"github.com/UpCloudLtd/upcloud-go-api/credentials"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
)

const (
	EnvDebugAPIBaseURL            = "UPCLOUD_DEBUG_API_BASE_URL"
	EnvDebugSkipCertificateVerify = "UPCLOUD_DEBUG_SKIP_CERTIFICATE_VERIFY"
	defaultUserAgentProduct       = "upcloud-go-api/v9"
)

// New creates a client with an API token.
func New(token string, opts ...ClientOption) (*ClientWithResponses, error) {
	provider, err := securityprovider.NewSecurityProviderBearerToken(token)
	if err != nil {
		return nil, err
	}
	opt := ClientOption(func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, provider.Intercept)
		return nil
	})
	return NewClientWithResponses(ServerUrlHttpsapiUpcloudCom, append([]ClientOption{WithDefaultUserAgent(), opt}, opts...)...)
}

// NewFromEnv creates a client from env or keyring. Set UPCLOUD_TOKEN or UPCLOUD_USERNAME/UPCLOUD_PASSWORD.
func NewFromEnv(opts ...ClientOption) (*ClientWithResponses, error) {
	creds, err := credentials.Parse(credentials.Credentials{})
	if err != nil {
		return nil, err
	}
	hasToken := creds.Token != ""
	hasBasic := creds.Username != "" && creds.Password != ""
	if hasToken && hasBasic {
		return nil, errors.New("only one authentication method (token or basic auth) can be provided")
	}
	if !hasToken && !hasBasic {
		return nil, errors.New("authentication credentials must be provided via environment variables or system keyring")
	}

	allOpts := []ClientOption{WithDefaultUserAgent(), WithCredentials(creds)}
	if os.Getenv(EnvDebugSkipCertificateVerify) == "1" {
		allOpts = append(allOpts, WithInsecureSkipVerify())
	}

	server := ServerUrlHttpsapiUpcloudCom
	if u := os.Getenv(EnvDebugAPIBaseURL); u != "" {
		server = u
	}
	return NewClientWithResponses(server, append(allOpts, opts...)...)
}

func defaultUserAgent() string {
	return defaultUserAgentProduct + " openapi/" + specVersion
}

// WithDefaultUserAgent injects the SDK User-Agent when one is not already set.
func WithDefaultUserAgent() ClientOption {
	return WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		if req.Header.Get("User-Agent") == "" {
			req.Header.Set("User-Agent", defaultUserAgent())
		}
		return nil
	})
}

func WithCredentials(creds credentials.Credentials) ClientOption {
	return func(c *Client) error {
		switch creds.Type() {
		case credentials.CredentialsTypeToken:
			provider, err := securityprovider.NewSecurityProviderBearerToken(creds.Token)
			if err != nil {
				return err
			}
			c.RequestEditors = append(c.RequestEditors, provider.Intercept)
			return nil
		case credentials.CredentialsTypeBasic:
			provider, err := securityprovider.NewSecurityProviderBasicAuth(creds.Username, creds.Password)
			if err != nil {
				return err
			}
			c.RequestEditors = append(c.RequestEditors, provider.Intercept)
			return nil
		default:
			return errors.New("credentials not defined (set UPCLOUD_TOKEN or UPCLOUD_USERNAME/UPCLOUD_PASSWORD)")
		}
	}
}

// WithInsecureSkipVerify skips TLS certificate verification. For debugging only.
func WithInsecureSkipVerify() ClientOption {
	return WithHTTPClient(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
		},
	})
}
