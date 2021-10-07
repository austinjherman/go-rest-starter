package logging

import (
	"aherman/src/http/response"
	"encoding/json"
	"errors"
)

// ServerError represents an extended ErrClientSafe
// in which we will bind the server error and log it out.
// This should not be used as an http response.
type ServerError struct {
	ClientErr response.ErrClientSafe
	ErrDetails error `json:"error_details"`
}

// BindClientErr todo
func (err *ServerError) BindClientErr(incoming *response.ErrClientSafe) *ServerError {
	if incoming == nil {
		return nil
	}
	(*err).ClientErr.Err = incoming.Err
	(*err).ClientErr.ErrDescription = incoming.ErrDescription
	(*err).ClientErr.ErrValidation = incoming.ErrValidation
	(*err).ClientErr.OK = incoming.OK
	(*err).ClientErr.Status = incoming.Status
	return err
}

// Error returns a string so Error implements the error interface
func (err *ServerError) Error() string {
	b, _ := json.Marshal(*err)
	return string(b)
}

// MarshalJSON todo
func (err ServerError) MarshalJSON() ([]byte, error) {
	anon := struct{
		ClientErr response.ErrClientSafe `json:"err_client"`
		ErrDetails string `json:"err_server"`
	}{
		ClientErr: err.ClientErr,
		ErrDetails: err.ErrDetails.Error(),
	}
	return json.Marshal(anon)
}

// NewServerError todo
func NewServerError(incoming error) *ServerError {
	serverErr := &ServerError{}
	clientSafeErr := &response.ErrClientSafe{}

	if errors.As(incoming, clientSafeErr) {
		(*serverErr).ErrDetails = errors.New("client safe")
		return serverErr
	}

	(*serverErr).ErrDetails = incoming
	return serverErr
}