package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// ErrClientSafe represents an error that is safe to return to the client.
// We shouldn't use this to expose any database or application-level errors;
// instead, this is intended to be used for providing information to the client.
// ErrClientSafe implements the Error interface.
type ErrClientSafe struct {
	Err            string `json:"error"`
	ErrDescription string `json:"error_description"`
	ErrValidation  []ValidationError `json:"error_validation,omitempty"`
	OK             bool   `json:"ok"`
	Status         int    `json:"status"`
}

// Error returns a string so ErrClientSafe implements the Error interface.
func (msg ErrClientSafe) Error() string {
	b, _ := json.Marshal(msg)
	return string(b)
}

// FormatValidationErrs will type-check the incoming error to see if it's 
// a validation error. If it is, it will format them nicely in msg.ErrValidation.
// If it's not, it won't do anything.
func (msg *ErrClientSafe) FormatValidationErrs(incoming error) {

	// check if the incoming error is of the validation error type.
	// If it is, we can give the user some useful information.
	var valErrs *validator.ValidationErrors = &validator.ValidationErrors{}
	if errors.As(incoming, valErrs) {
		msg.ErrValidation = NewJSONFormatter().Descriptive(valErrs)
	}
}

// Error checks to see if the provided error is of the ErrorMsg type.
// If it is, then we can safely return the error to the client. If it's 
// not, this is more likely a deeper application issue that we'll mask 
// by returning ErrUnknown.
func Error(err error) (int, error) {
	if apiError, ok := err.(ErrClientSafe); ok {
		return apiError.Status, apiError
	}
	return ErrUnknown.Status, ErrUnknown
}

// ErrBadRequest is the HTTP error to be returned when the client
// makes a request that doesn't pass validation
var ErrBadRequest ErrClientSafe = ErrClientSafe{
	Err: "bad_request",
	ErrDescription: "A validation error occurred",
	OK:          false,
	Status:      http.StatusBadRequest,
}

// ErrClientUnknown is the HTTP status code to be returned when 
// requests are made from an unknown client.
var ErrClientUnknown ErrClientSafe = ErrClientSafe{
	Err: "client_unknown",
	ErrDescription: "Uknown client",
	OK: false,
	Status: http.StatusUnauthorized,
}

// ErrInvitationCodeInvalid is the HTTP error to be returned if a 
// user provides an invalid invitation code
var ErrInvitationCodeInvalid ErrClientSafe = ErrClientSafe{
	Err: "invitation_code_invalid",
	ErrDescription: "Invalid invitation code",
	OK: false,
	Status: http.StatusUnauthorized,
}

// ErrNotImplemented is the HTTP error to be returned when the
// requested resource is not implemented
var ErrNotImplemented ErrClientSafe = ErrClientSafe{
	Err: "not_implemented",
	ErrDescription: "Not implemented",
	OK: false,
	Status: http.StatusNotImplemented,
}

// ErrResourceNotFound is the HTTP error to be returned when the
// requested resource is not implemented
var ErrResourceNotFound ErrClientSafe = ErrClientSafe{
	Err: "resource_not_found",
	ErrDescription: "Resource not found",
	OK: false,
	Status: http.StatusNotFound,
}

// ErrResourceNotUpdated is the HTTP error to be returned when the
// requested resource is not implemented
var ErrResourceNotUpdated ErrClientSafe = ErrClientSafe{
	Err: "resource_not_updated",
	ErrDescription: "Resource couldn't be updated",
	OK: false,
	Status: http.StatusInternalServerError,
}

// ErrTokenExpired is the HTTP error to be returned when
// an access token is expired
var ErrTokenExpired ErrClientSafe = ErrClientSafe{
	Err: "token_expired",
	ErrDescription: "Expired token",
	OK: false,
	Status: http.StatusUnauthorized,
}

// ErrTokenInvalid is the HTTP error to be returned when a token is invalid
var ErrTokenInvalid ErrClientSafe = ErrClientSafe{
	Err: "token_invalid",
	ErrDescription: "Invalid token",
	OK: false,
	Status: http.StatusUnauthorized,
}

// ErrTokenNotFound TODO
var ErrTokenNotFound ErrClientSafe = ErrClientSafe{
	Err: "token_not_found",
	ErrDescription: "Token not found",
	OK: false,
	Status: http.StatusNotFound,
}

// ErrTokenTypeInvalid is the HTTP error to be returned when
// a refresh token was sent to the access token route
var ErrTokenTypeInvalid ErrClientSafe = ErrClientSafe{
	Err: "token_type_invalid",
	ErrDescription: "Invalid token type",
	OK: false,
	Status: http.StatusBadRequest,
}

// ErrUnknown TODO
var ErrUnknown ErrClientSafe = ErrClientSafe{
	Err: "err_unknown",
	ErrDescription: "An unknown error occured",
	OK: false,
	Status: http.StatusInternalServerError,
}

// ErrUpdate is the HTTP error to be returned when the
// user gives an invalid password
var ErrUpdate ErrClientSafe = ErrClientSafe{
	Err: "update",
	ErrDescription: "Could not update resource",
	OK: false,
	Status: http.StatusInternalServerError,
}

// ErrUserEmailAlreadyRegistered is the HTTP error to be returned when
// the user tries to sign up with an email that is in use
var ErrUserEmailAlreadyRegistered ErrClientSafe = ErrClientSafe{
	Err: "user_email_already_registered",
	ErrDescription: "Email address already registered",
	OK: false,
	Status: http.StatusBadRequest,
}

// ErrUserEmailNotFound is the HTTP error to be returned when the
// email could not be found
var ErrUserEmailNotFound ErrClientSafe = ErrClientSafe{
	Err: "user_email_not_found",
	ErrDescription: "Email not found",
	OK: false,
	Status: http.StatusNotFound,
}

// ErrUserLogout is the HTTP error to be returned if there was an
// issue logging the user out
var ErrUserLogout ErrClientSafe = ErrClientSafe{
	Err: "user_logout",
	ErrDescription: "Could not log out",
	OK: false,
	Status: http.StatusInternalServerError,
}

// ErrUserPasswordInvalid is the HTTP error to be returned when the
// user gives an invalid password
var ErrUserPasswordInvalid ErrClientSafe = ErrClientSafe{
	Err: "user_password_invalid",
	ErrDescription: "Invalid password",
	OK: false,
	Status: http.StatusUnauthorized,
}

// ErrUserNotFound is the HTTP error to be returned when
// the user can't be found
var ErrUserNotFound ErrClientSafe = ErrClientSafe{
	Err: "user_not_found",
	ErrDescription: "User not found",
	OK: false,
	Status: http.StatusNotFound,
}

// ErrValidation is the HTTP error to be returned when
// validation fails
var ErrValidation ErrClientSafe = ErrClientSafe{
	Err: "err_validation",
	ErrDescription: "Validation error",
	OK: false,
	Status: http.StatusBadRequest,
}