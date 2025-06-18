package credentials_test

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/credentials"
	"github.com/stretchr/testify/assert"
	"github.com/zalando/go-keyring"
)

func ensureTokenInKeyring(t *testing.T) string {
	t.Helper()

	token, _ := keyring.Get("UpCloud", "")
	if token == "" {
		token = "unittest_token"
		err := keyring.Set("UpCloud", "", token)
		assert.NoError(t, err)
		t.Cleanup(func() {
			assert.NoError(t, keyring.Delete("UpCloud", ""))
		})
	}

	return token
}

func TestParse_PasswordFromKeyring(t *testing.T) {
	ensureTokenInKeyring(t)
	sources := []string{"configuration", "environment"}

	for _, source := range sources {
		t.Run(source, func(t *testing.T) {
			cfg := credentials.Credentials{}
			if source == "configuration" {
				cfg.Username = "unittest"
			} else {
				t.Setenv("UPCLOUD_USERNAME", "unittest")
			}

			err := keyring.Set("UpCloud", "unittest", "unittest_password")
			assert.NoError(t, err)
			t.Cleanup(func() {
				assert.NoError(t, keyring.Delete("UpCloud", "unittest"))
			})

			creds, err := credentials.Parse(cfg)
			assert.NoError(t, err)

			assert.Equal(t, "unittest", creds.Username)
			assert.Equal(t, "unittest_password", creds.Password)
			assert.Equal(t, "keyring", string(creds.Source()))
			assert.Equal(t, "basic", string(creds.Type()))
		})
	}
}

func TestParse_TokenFromKeyring(t *testing.T) {
	token := ensureTokenInKeyring(t)

	creds, err := credentials.Parse(credentials.Credentials{})
	assert.NoError(t, err)

	assert.Equal(t, "", creds.Username)
	assert.Equal(t, "", creds.Password)
	assert.Equal(t, token, creds.Token)
	assert.Equal(t, "keyring", string(creds.Source()))
	assert.Equal(t, "token", string(creds.Type()))
}

func TestParse_BasicFromEnv(t *testing.T) {
	t.Setenv("UPCLOUD_USERNAME", "unittest")
	t.Setenv("UPCLOUD_PASSWORD", "unittest_password")
	t.Setenv("UPCLOUD_TOKEN", "")

	creds, err := credentials.Parse(credentials.Credentials{})
	assert.NoError(t, err)

	assert.Equal(t, "unittest", creds.Username)
	assert.Equal(t, "unittest_password", creds.Password)
	assert.Equal(t, "", creds.Token)
	assert.Equal(t, "environment", string(creds.Source()))
	assert.Equal(t, "basic", string(creds.Type()))
}

func TestParse_TokenFromEnv(t *testing.T) {
	t.Setenv("UPCLOUD_USERNAME", "unittest")
	t.Setenv("UPCLOUD_PASSWORD", "unittest_password")
	t.Setenv("UPCLOUD_TOKEN", "unittest_token")

	creds, err := credentials.Parse(credentials.Credentials{})
	assert.NoError(t, err)

	assert.Equal(t, "", creds.Username)
	assert.Equal(t, "", creds.Password)
	assert.Equal(t, "unittest_token", creds.Token)
	assert.Equal(t, "environment", string(creds.Source()))
	assert.Equal(t, "token", string(creds.Type()))
}

func TestConfig_OverrideWithParameters(t *testing.T) {
	t.Setenv("UPCLOUD_USERNAME", "")
	t.Setenv("UPCLOUD_PASSWORD", "")
	t.Setenv("UPCLOUD_TOKEN", "unittest_token")

	creds, err := credentials.Parse(credentials.Credentials{
		Username: "override_user",
		Password: "override_pass",
	})
	assert.NoError(t, err)

	assert.Equal(t, "override_user", creds.Username)
	assert.Equal(t, "override_pass", creds.Password)
	assert.Equal(t, "", creds.Token)
	assert.Equal(t, "configuration", string(creds.Source()))
	assert.Equal(t, "basic", string(creds.Type()))
}
