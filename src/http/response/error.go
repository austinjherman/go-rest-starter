package response

import (
	"net/http"
)

// ClientSafeErrorMsg represents an error that is safe to return to the client.
// We shouldn't use this to expose any database or application-level errors;
// instead, this is intended to be used for providing information to the client.
// ErrorMsg implements the Error interface.
type ClientSafeErrorMsg struct {
	OK          bool   `json:"ok"`
	Status      int    `json:"status"`
	Description string `json:"description"`
}

// Error returns a string so ErrorMsg implements the Error interface.
func (err ClientSafeErrorMsg) Error() string {
	return err.Description
}

// Error checks to see if the provided error is of the ErrorMsg type.
// If it is, then we can safely return the error to the client. If it's 
// not, this is more likely a deeper application issue that we'll mask 
// by returning ErrUnknown.
func Error(err error) (int, ClientSafeErrorMsg) {
	if apiError, ok := err.(ClientSafeErrorMsg); ok {
		return apiError.Status, apiError
	}
	
	return ErrUnknown.Status, ErrUnknown
}

// ErrAccessTokenInvalid TODO
var ErrAccessTokenInvalid ClientSafeErrorMsg = ClientSafeErrorMsg{
	OK:          false,
	Status:      http.StatusUnauthorized,
	Description: "Invalid access token",
}

// ErrAccessTokenExpired TODO
var ErrAccessTokenExpired ClientSafeErrorMsg = ClientSafeErrorMsg{
	OK:          false,
	Status:      http.StatusUnauthorized,
	Description: "Expired access token",
}

// ErrAccessTokenParse TODO
var ErrAccessTokenParse ClientSafeErrorMsg = ClientSafeErrorMsg{
	OK:          false,
	Status:      http.StatusInternalServerError,
	Description: "Error parsing token",
}

// ErrBadCredentials TODO
var ErrBadCredentials ClientSafeErrorMsg = ClientSafeErrorMsg{
	OK:          false,
	Status:      http.StatusUnauthorized,
	Description: "Invalid credentials",
}

// ErrBadRequest TODO
var ErrBadRequest ClientSafeErrorMsg = ClientSafeErrorMsg{
	OK:          false,
	Status:      http.StatusBadRequest,
	Description: "Bad request",
}

// ErrEmailNotAvailable TODO
var ErrEmailNotAvailable ClientSafeErrorMsg = ClientSafeErrorMsg{
	OK:          false,
	Status:      http.StatusBadRequest,
	Description: "Email address already registered",
}

// ErrEmailNotFound TODO
var ErrEmailNotFound ClientSafeErrorMsg = ClientSafeErrorMsg{
	OK:          false,
	Status:      http.StatusNotFound,
	Description: "Email address not found",
}

// ErrInvalidPassword TODO
var ErrInvalidPassword ClientSafeErrorMsg = ClientSafeErrorMsg{
	OK:          false,
	Status:      http.StatusUnauthorized,
	Description: "Invalid password",
}

// ErrLogout TODO
var ErrLogout ClientSafeErrorMsg = ClientSafeErrorMsg{
	OK:          false,
	Status:      http.StatusInternalServerError,
	Description: "Could not log out",
}

// ErrNoInvitation TODO
var ErrNoInvitation ClientSafeErrorMsg = ClientSafeErrorMsg{
	OK:          false,
	Status:      http.StatusUnauthorized,
	Description: "Invitation code required",
}

// ErrNotATrustedClient TODO
var ErrNotATrustedClient ClientSafeErrorMsg = ClientSafeErrorMsg{
	OK:          false,
	Status:      http.StatusUnauthorized,
	Description: "Only trusted clients can request this resource",
}

// ErrNotFound TODO
var ErrNotFound ClientSafeErrorMsg = ClientSafeErrorMsg{
	OK:          false,
	Status:      http.StatusNotFound,
	Description: "Not found",
}

// ErrNotImplemented TODO
var ErrNotImplemented ClientSafeErrorMsg = ClientSafeErrorMsg{
	OK:          false,
	Status:      http.StatusNotImplemented,
	Description: "Not implemented",
}

// ErrValidation TODO
var ErrValidation ClientSafeErrorMsg = ClientSafeErrorMsg{
	OK:          false,
	Status:      http.StatusBadRequest,
	Description: "A validation error occurred",
}

// ErrUnauthorized TODO
var ErrUnauthorized ClientSafeErrorMsg = ClientSafeErrorMsg{
	OK:          false,
	Status:      http.StatusUnauthorized,
	Description: "Unauthorized",
}

// ErrUnknown TODO
var ErrUnknown ClientSafeErrorMsg = ClientSafeErrorMsg{
	OK:          false,
	Status:      http.StatusInternalServerError,
	Description: "An unknown error occured",
}
