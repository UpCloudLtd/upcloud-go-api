package credentials

type CredentialsSource string

const (
	CredentialsSourceConfiguration CredentialsSource = "configuration"
	CredentialsSourceEnvironment   CredentialsSource = "environment"
	CredentialsSourceKeyring       CredentialsSource = "keyring"
)

type CredentialsType string

const (
	CredentialsTypeBasic CredentialsType = "basic"
	CredentialsTypeToken CredentialsType = "token"
)

type Credentials struct {
	Username string
	Password string
	Token    string

	source CredentialsSource
}

func (c Credentials) IsDefined() bool {
	return (c.Username != "" && c.Password != "") || c.Token != ""
}

func (c Credentials) Source() CredentialsSource {
	return c.source
}

func (c Credentials) Type() CredentialsType {
	if c.Token != "" {
		return CredentialsTypeToken
	}
	if c.Username != "" && c.Password != "" {
		return CredentialsTypeBasic
	}
	return ""
}

func (c Credentials) clean(source CredentialsSource) Credentials {
	c.source = source
	if c.Token != "" {
		c.Username = ""
		c.Password = ""
	}

	return c
}
