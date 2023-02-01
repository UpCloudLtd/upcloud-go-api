package upcloud

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNotFoundError(t *testing.T) {
	notFoundErr := &Error{
		ErrorCode: ErrCodeServerNotFound,
		Status:    http.StatusNotFound,
	}

	notFoundErr2 := &Error{
		ErrorCode: ErrCodeRouterNotFound,
		Status:    http.StatusNotFound,
	}

	notFoundProblem := &Problem{
		Status: http.StatusNotFound,
	}

	assert.True(t, IsNotFoundError(notFoundErr))
	assert.True(t, IsNotFoundError(notFoundErr2))
	assert.True(t, IsNotFoundError(notFoundProblem))

	otherError := &Error{
		ErrorCode: ErrCodeDBExists,
		Status:    http.StatusConflict,
	}

	otherProblem := &Problem{
		Status: http.StatusBadRequest,
	}

	assert.False(t, IsNotFoundError(otherError))
	assert.False(t, IsNotFoundError(otherProblem))
}

func TestIsAlreadyExistsError(t *testing.T) {
	alreadyExistsErr := &Error{
		ErrorCode: ErrCodeDBExists,
		Status:    http.StatusConflict,
	}

	alreadyExistsErr2 := &Error{
		ErrorCode: ErrCodeTagExists,
		Status:    http.StatusConflict,
	}

	alreadyExistsProblem := &Problem{
		Type:   "https://developers.upcloud.com/1.3/errors#ERROR_RESOURCE_ALREADY_EXISTS",
		Status: http.StatusBadRequest,
	}

	assert.True(t, IsAlreadyExistsError(alreadyExistsErr))
	assert.True(t, IsAlreadyExistsError(alreadyExistsErr2))
	assert.True(t, IsAlreadyExistsError(alreadyExistsProblem))

	otherError := &Error{
		ErrorCode: ErrCodeAuthenticationFailed,
		Status:    http.StatusConflict,
	}

	otherProblem := &Problem{
		Status: http.StatusBadRequest,
	}

	assert.False(t, IsAlreadyExistsError(otherError))
	assert.False(t, IsAlreadyExistsError(otherProblem))
}

func TestIsAuthenticationFailedError(t *testing.T) {
	authFailedErr := &Error{
		ErrorCode: ErrCodeAuthenticationFailed,
		Status:    http.StatusUnauthorized,
	}

	authFailedProblem := &Problem{
		Type:   "https://developers.upcloud.com/1.3/errors#ERROR_AUTHENTICATION_FAILED",
		Status: http.StatusUnauthorized,
	}

	assert.True(t, IsAuthenticationFailedError(authFailedErr))
	assert.True(t, IsAuthenticationFailedError(authFailedProblem))

	otherError := &Error{
		ErrorCode: ErrCodeAuthenticationFailed,
		Status:    http.StatusConflict,
	}

	otherProblem := &Problem{
		Status: http.StatusBadRequest,
	}

	assert.False(t, IsAuthenticationFailedError(otherError))
	assert.False(t, IsAuthenticationFailedError(otherProblem))
}
