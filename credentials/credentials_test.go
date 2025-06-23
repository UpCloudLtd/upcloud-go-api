package credentials_test

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/credentials"
	"github.com/stretchr/testify/assert"
)

func TestCredentials_IsDefined(t *testing.T) {
	testdata := []struct {
		name     string
		creds    credentials.Credentials
		expected bool
	}{
		{
			name:     "empty credentials",
			creds:    credentials.Credentials{},
			expected: false,
		},
		{
			name:     "only username",
			creds:    credentials.Credentials{Username: "user"},
			expected: false,
		},
		{
			name:     "only password",
			creds:    credentials.Credentials{Password: "pass"},
			expected: false,
		},
		{
			name: "only token",
			creds: credentials.Credentials{
				Token: "token",
			},
			expected: true,
		},
		{
			name:     "username and password",
			creds:    credentials.Credentials{Username: "user", Password: "pass"},
			expected: true,
		},
		{
			name:     "username, password and token",
			creds:    credentials.Credentials{Username: "user", Password: "pass", Token: "token"},
			expected: true,
		},
	}

	for _, test := range testdata {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.creds.IsDefined(), test.expected)
		})
	}
}
