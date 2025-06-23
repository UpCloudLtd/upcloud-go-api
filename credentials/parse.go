package credentials

import (
	"fmt"
	"os"

	"github.com/zalando/go-keyring"
)

// KeyringServiceName is the name of the service to use when using the system keyring
const KeyringServiceName string = "UpCloud"

// KeyringTokenUser is the username to use for the token in the system keyring
const KeyringTokenUser string = ""

func readFromEnv(creds Credentials) Credentials {
	if creds.Username == "" {
		creds.Username = os.Getenv("UPCLOUD_USERNAME")
	}
	if creds.Password == "" {
		creds.Password = os.Getenv("UPCLOUD_PASSWORD")
	}
	if creds.Token == "" {
		creds.Token = os.Getenv("UPCLOUD_TOKEN")
	}

	return creds
}

func readFromKeyring(creds Credentials) Credentials {
	if creds.Username != "" {
		creds := Credentials{
			Username: creds.Username,
		}

		password, _ := keyring.Get(KeyringServiceName, creds.Username)
		if password != "" {
			creds.Password = password
			return creds
		}
	}

	token, _ := keyring.Get(KeyringServiceName, KeyringTokenUser)

	return Credentials{
		Token: token,
	}
}

// Parse reads credentials from environment variables or the system keyring and allows overriding these with parameters. If both basic auth and token are provided, basic auth credentials will be omitted from the return value. Any credentials defined in the parameters will take precedence over those read from the environment and any credentials defined in the environment will take precedence over those read from the keyring.
func Parse(config Credentials) (Credentials, error) {
	creds := config
	if creds.IsDefined() {
		return creds.clean(CredentialsSourceConfiguration), nil
	}

	creds = readFromEnv(creds)
	if creds.IsDefined() {
		return creds.clean(CredentialsSourceEnvironment), nil
	}

	creds = readFromKeyring(creds)
	if creds.IsDefined() {
		return creds.clean(CredentialsSourceKeyring), nil
	}

	return Credentials{}, fmt.Errorf("credentials not found, these must be set in configuration, via environment variables or in the system keyring (%s)", KeyringServiceName)
}
