package upcloud

import "encoding/json"

// Account represents an account
type Account struct {
	Credits  float64 `xml:"credits"`
	UserName string  `xml:"username"`
}

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (s *Account) UnmarshalJSON(b []byte) error {
	type localAccount Account

	v := struct {
		Account localAccount `json:"account"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	(*s) = Account(v.Account)

	return nil
}
