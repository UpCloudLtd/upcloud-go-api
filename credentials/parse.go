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

func readFromEnv() Credentials {
	return Credentials{
		Username: os.Getenv("UPCLOUD_USERNAME"),
		Password: os.Getenv("UPCLOUD_PASSWORD"),
		Token:    os.Getenv("UPCLOUD_TOKEN"),
	}
}

func readFromKeyring(username string) Credentials {
	token, err := keyring.Get(KeyringServiceName, KeyringTokenUser)
	fmt.Printf("Failed to read from keyring: %s", err)
	if err == nil {
		return Credentials{
			Token: token,
		}
	}

	creds := Credentials{
		Username: username,
	}

	password, err := keyring.Get(KeyringServiceName, username)
	if err == nil {
		creds.Password = password
	}

	return creds
}

// Parse reads credentials from environment variables or the system keyring and allows overriding these with parameters. If both basic auth and token are provided, basic auth credentials will be omitted from the return value. Any credentials defined in the parameters will take precedence over those read from the environment and any credentials defined in the environment will take precedence over those read from the keyring.
func Parse(config Credentials) (Credentials, error) {
	creds := config
	if creds.IsDefined() {
		return creds.clean(CredentialsSourceConfiguration), nil
	}

	creds = readFromEnv()
	if creds.IsDefined() {
		return creds.clean(CredentialsSourceEnvironment), nil
	}

	creds = readFromKeyring(creds.Username)
	if creds.IsDefined() {
		return creds.clean(CredentialsSourceKeyring), nil
	}

	return Credentials{}, fmt.Errorf("credentials not found, these must be set in configuration, via environment variables or in the system keyring (%s)", KeyringServiceName)
}
