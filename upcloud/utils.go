package upcloud

import (
	"encoding/json"
)

// Constants
const (
	True  Boolean = 1
	False Boolean = -1
	Empty Boolean = 0
)

// Boolean is a custom boolean type that allows
// for custom marshalling and unmarshalling and
// for an empty value that isn't false so we can
// distinguish between true, false and not set.
type Boolean int

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (b *Boolean) UnmarshalJSON(buf []byte) error {
	str := string(buf)
	if str == `true` ||
		str == `"true"` ||
		str == `"yes"` ||
		str == `1` ||
		str == `"1"` {
		(*b) = 1
		return nil
	}

	(*b) = -1
	return nil
}

// MarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (b *Boolean) MarshalJSON() ([]byte, error) {
	if (*b) == 1 {
		return []byte(`"yes"`), nil
	}

	return []byte(`"no"`), nil
}

// Bool converts to a standard bool value
func (b *Boolean) Bool() bool {
	return (*b) >= True
}

// Empty checks if this bool is empty
func (b *Boolean) Empty() bool {
	return (*b) == Empty
}

// String returns a string representation
func (b *Boolean) String() string {
	if (*b) >= True {
		return "True"
	}

	if (*b) <= False {
		return "False"
	}

	return "Empty"
}

// FromBool converts from a standard bool values
func FromBool(v bool) Boolean {
	if v {
		return True
	}
	return False
}

func StringPtr(v string) *string {
	return &v
}

func BoolPtr(v bool) *bool {
	return &v
}

func IntPtr(v int) *int {
	return &v
}

// ServerUUIDSlice is a slice of string.
// It exists to allow for a custom JSON unmarshaller.
type ServerUUIDSlice []string

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (s *ServerUUIDSlice) UnmarshalJSON(b []byte) error {
	v := struct {
		ServerUUIDs []string `json:"server"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	(*s) = v.ServerUUIDs

	return nil
}

// MarshalJSON is a custom marshaller that deals with deeply embedded values.
func (s *ServerUUIDSlice) MarshalJSON() ([]byte, error) {
	v := struct {
		ServerUUIDs []string `json:"server"`
	}{
		ServerUUIDs: *s,
	}
	return json.Marshal(&v)
}
