package response

import (
	"net/http"
)

// ErrorMsg represents the structure of errors as they occur in the application.
type ErrorMsg struct {
	OK          bool   `json:"ok"`
	Status      int    `json:"status"`
	Description string `json:"description"`
}

func (err ErrorMsg) Error() string {
	return err.Description
}

// Error TODO
func Error(err error) (int, ErrorMsg) {
	// This the best way to log?
	// trace := make([]byte, 1024)
	// runtime.Stack(trace, true)
	// log.Printf("ERROR: %s\n%s", err, trace)
	if apiError, ok := err.(ErrorMsg); ok {
		return apiError.Status, apiError
	}

	return ErrUnknown.Status, ErrUnknown
}

// ErrBadCredentials TODO
var ErrBadCredentials ErrorMsg = ErrorMsg{
	OK:          false,
	Status:      http.StatusUnauthorized,
	Description: "Your credentials are incorrect.",
}

// ErrBadRequest TODO
var ErrBadRequest ErrorMsg = ErrorMsg{
	OK:          false,
	Status:      http.StatusBadRequest,
	Description: "You've made a bad request.",
}

// ErrDatabase TODO
var ErrDatabase ErrorMsg = ErrorMsg{
	OK:          false,
	Status:      http.StatusInternalServerError,
	Description: "An internal server error occured.",
}

// ErrEmailNotAvailable TODO
var ErrEmailNotAvailable ErrorMsg = ErrorMsg{
	OK:          false,
	Status:      http.StatusBadRequest,
	Description: "That email address is already registered.",
}

// ErrEmailNotFound TODO
var ErrEmailNotFound ErrorMsg = ErrorMsg{
	OK:          false,
	Status:      http.StatusNotFound,
	Description: "The provided email address could not be found.",
}

// ErrNoInvitation TODO
var ErrNoInvitation ErrorMsg = ErrorMsg{
	OK:          false,
	Status:      http.StatusNotFound,
	Description: "You need an invitation code to use this application.",
}

// ErrNotFound TODO
var ErrNotFound ErrorMsg = ErrorMsg{
	OK:          false,
	Status:      http.StatusNotFound,
	Description: "The requested resource could not be found.",
}

// ErrUserLoggedIn TODO
var ErrUserLoggedIn ErrorMsg = ErrorMsg{
	OK:          false,
	Status:      http.StatusBadRequest,
	Description: "This user is already logged in.",
}

// ErrValidation TODO
var ErrValidation ErrorMsg = ErrorMsg{
	OK:          false,
	Status:      http.StatusBadRequest,
	Description: "A validation error occurred.",
}

// ErrUnauthorized TODO
var ErrUnauthorized ErrorMsg = ErrorMsg{
	OK:          false,
	Status:      http.StatusUnauthorized,
	Description: "You need to be logged in.",
}

// ErrUnknown TODO
var ErrUnknown ErrorMsg = ErrorMsg{
	OK:          false,
	Status:      http.StatusInternalServerError,
	Description: "An unknown error occured.",
}

// ErrUsernameNotAvailable TODO
var ErrUsernameNotAvailable ErrorMsg = ErrorMsg{
	OK:          false,
	Status:      http.StatusBadRequest,
	Description: "That username is already in use.",
}
