package errs_test

import (
	"errors"
	"testing"

	"github.com/avegner/utils/errs"
)

var (
	errNil    error
	errNonNil = errors.New("initial error")
	errCb     = errors.New("last error")
)

func TestFinallyPropagatesCbError(t *testing.T) {
	testFinally(t, errNil, errNil, errNil)
	testFinally(t, errNil, errCb, errCb)
}

func TestFinallySavesFuncError(t *testing.T) {
	testFinally(t, errNonNil, errNil, errNonNil)
	testFinally(t, errNonNil, errCb, errNonNil)
}

func testFinally(t *testing.T, initErr, cbErr, wantErr error) {
	err := initErr
	defer func() {
		if err != wantErr {
			t.Fatalf("got '%v' error, want '%v' error", err, wantErr)
		}
	}()
	defer errs.Finally(func() error {
		return cbErr
	}, &err)
}
