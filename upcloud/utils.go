package upcloud

// Boolean is a custom boolean type that allows
// for custom marshalling and unmarshalling
type Boolean bool

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (b *Boolean) UnmarshalJSON(buf []byte) error {
	str := string(buf)
	if str == `true` ||
		str == `"true"` ||
		str == `"yes"` ||
		str == `1` ||
		str == `"1"` {
		(*b) = true
		return nil
	}

	(*b) = false
	return nil
}

// MarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (b Boolean) MarshalJSON() ([]byte, error) {
	if b {
		return []byte(`"yes"`), nil
	}

	return []byte(`"no"`), nil
}
