/**************************************************************************************
  Error Facade
	-------------
	The error fade is useful because todo
**************************************************************************************/

package error

import (
	"aherman/src/http/response"
	"aherman/src/logging"
	"errors"
	"io"
)

// Compiler todo
type Compiler interface {
	ShouldContinue(error) (bool, *response.ErrClientSafe)
}

// Error todo
type Error struct {
	ServerWriter io.Writer
}

/**************************************************************************************
  Public Functions
**************************************************************************************/

// ShouldContinue takes an error, incoming. If incoming is nil, ShouldContinue will
// indicate to the caller that it can proceed. If the incoming error is a client-safe
// error, ShouldContinue will return the original error with formatted validation
// messages if there are any. It the error is of any other type, ShouldContinue will
// create a server error and write it to the Error stream.
func (facade *Error) ShouldContinue(
	incoming error,
	defaultErr *response.ErrClientSafe,
) (bool, *response.ErrClientSafe) {

	if incoming == nil {
		return true, nil
	}

	clientSafeResponse := &response.ErrClientSafe{}
	if errors.As(incoming, &clientSafeResponse) {
		clientSafeResponse.FormatValidationErrs(incoming)
		return false, clientSafeResponse
	}

	serverErr := logging.NewServerError(incoming).BindClientErr(defaultErr)
	facade.ServerWriter.Write([]byte(serverErr.Error()))

	return false, defaultErr
}
