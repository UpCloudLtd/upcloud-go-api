package upcloud

type Boolean bool

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

func (b Boolean) MarshalJSON() ([]byte, error) {
	if b {
		return []byte(`"yes"`), nil
	}

	return []byte(`"no"`), nil
}
