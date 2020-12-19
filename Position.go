package geolocation

import (
	"syscall/js"
	"time"
)

// Position is a container for the geolocation information.
type Position struct {
	Coords *Coords

	// Timestamp is the acquisition time of the position.
	Timestamp time.Time
}

// CurrentPosition requests the location of the device.
func CurrentPosition(options *PositionOptions) (*Position, *Error) {
	var errChan = make(chan *Error, 1)

	var position *Position

	successCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		result := args[0]

		resultTS := result.Get("timestamp").Int()

		position = &Position{
			Coords:    toCoord(result.Get("coords")),
			Timestamp: time.Unix(int64(resultTS)/1000, int64(resultTS)%1000),
		}

		errChan <- nil

		return nil
	})
	defer successCallback.Release()

	errorCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		err := toError(args[0])
		errChan <- &err

		return nil
	})
	defer errorCallback.Release()

	geolocation().Call("getCurrentPosition", successCallback, errorCallback, options.jsValue())

	return position, <-errChan
}
