package upcloud

import "fmt"

/**
Represents an error
*/
type Error struct {
	ErrorCode    string `xml:"error_code"`
	ErrorMessage string `xml:"error_message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s (%s)", e.ErrorMessage, e.ErrorCode)
}
