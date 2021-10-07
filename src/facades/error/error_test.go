package error

import (
	"aherman/src/http/response"
	"bytes"
	"errors"
	"fmt"
	"testing"
)

var (
	writer *bytes.Buffer = &bytes.Buffer{}
	errorFacade *Error = &Error{ ServerWriter: writer }
)

func TestShouldContinueIfErrIsNil(t *testing.T) {
	ok, _ := errorFacade.ShouldContinue(nil, nil)
	if !ok {
		t.Fatalf("ok should be true")
	}
}

func TestShouldContinueErrClientSafe(t *testing.T) {
	err := &response.ErrClientSafe{ Err: "test_client_error" }
	ok, httpResponse := errorFacade.ShouldContinue(err, nil)
	if ok {
		t.Fatalf("ok should be false")
	}
	if !errors.As(*httpResponse, &response.ErrClientSafe{}) {
		t.Fatal("response should be client safe")
	}
	if err.Err != "test_client_error" {
		t.Fatalf("the error code should be \"test_client_error\"")
	}
}

func TestShouldContinueErrServer(t *testing.T) {
	appErr := errors.New("test_server_error")
	defaultErr := &response.ErrClientSafe{ Err: "test_client_error" }

	ok, httpResponse := errorFacade.ShouldContinue(appErr, defaultErr)
	if ok {
		t.Fatalf("ok should be false")
	}

	fmt.Println(writer.String())

	if defaultErr != httpResponse {
		t.Fatalf("the http response should be equal to the default error")
	}
}