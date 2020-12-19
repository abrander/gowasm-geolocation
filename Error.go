package geolocation

import (
	"syscall/js"
)

// Error is an error type specific for this package. It satisfiers
// the error interface so you may pass it along as a regular error.
type Error string

// Error implements the error interface.
func (e Error) Error() string {
	return string(e)
}

// Temporary returns true if the error is temporary and the operation
// can be retried.
func (e Error) Temporary() bool {
	return e == ErrPositionUnavailable || e == ErrTimeout
}

const (
	// ErrPermissionDenied will be returned if the location
	// acquisition process failed because the document does
	// not have permission to use the Geolocation API.
	ErrPermissionDenied = Error("permission denied")

	// ErrPositionUnavailable will be returned if the position
	// of the device could not be determined. For instance, one
	// or more of the location providers used in the location
	// acquisition process reported an internal error that
	// caused the process to fail entirely.
	ErrPositionUnavailable = Error("position unavailable")

	// ErrTimeout will be returned if the length of time
	// specified by the Timeout property has elapsed before
	// the implementation could successfully acquire a new
	// GeolocationPosition object.
	ErrTimeout = Error("timeout")
)

func toError(value js.Value) Error {
	const PERMISSION_DENIED = 1
	const POSITION_UNAVAILABLE = 2
	const TIMEOUT = 3

	code := value.Get("code").Int()

	switch code {
	case PERMISSION_DENIED:
		return ErrPermissionDenied
	case POSITION_UNAVAILABLE:
		return ErrPositionUnavailable
	case TIMEOUT:
		return ErrTimeout
	}

	// We should never arrive here, but just in case we catch this in
	// case some browsers does not adhere to the spec.
	return Error(value.Get("message").String())
}
