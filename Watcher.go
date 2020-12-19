package geolocation

import (
	"syscall/js"
	"time"
)

// Watcher provides a convenient interface for watching the device position.
type Watcher struct {
	handle int

	errCallback js.Func
	posCallback js.Func

	errChan chan Error
	posChan chan *Position

	closed bool
}

const (
	// ErrWatcherClosed will be returned from Next if the watcher is closed.
	ErrWatcherClosed = Error("watcher closed")
)

// Next will return a position or an error. The error should be checked using
// Temporary(). If Temporary() returns true, the operation may be retried. If
// it returns false, you must consider the watcher closed.
func (w *Watcher) Next() (*Position, *Error) {
	if w.closed {
		err := ErrWatcherClosed

		return nil, &err
	}

	select {
	case pos := <-w.posChan:
		return pos, nil

	case err := <-w.errChan:
		return nil, &err
	}
}

// Chans returns two channels. One carrying location updates, and the
// other carrying errors. This is useful for getting updates in a
// select loop.
func (w Watcher) Chans() (chan *Position, chan Error) {
	return w.posChan, w.errChan
}

// Close unregisters the watcher with the host device.
func (w *Watcher) Close() {
	if w.closed {
		return
	}

	w.closed = true

	geolocation().Call("clearWatch", w.handle)

	w.errCallback.Release()
	w.posCallback.Release()

	close(w.errChan)
	close(w.posChan)
}

// WatchPosition will continually watch the device position.
func WatchPosition(options *PositionOptions) *Watcher {
	positionChan := make(chan *Position, 1)
	errChan := make(chan Error, 1)

	successCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		result := args[0]

		resultTS := result.Get("timestamp").Int()

		position := &Position{
			Coords:    toCoord(result.Get("coords")),
			Timestamp: time.Unix(int64(resultTS)/1000, int64(resultTS)%1000),
		}

		positionChan <- position

		return nil
	})

	errorCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		err := toError(args[0])
		errChan <- err

		return nil
	})

	jsHandle := geolocation().Call("watchPosition", successCallback, errorCallback, options.jsValue())
	handle := jsHandle.Int()

	w := &Watcher{
		handle: handle,

		errCallback: errorCallback,
		posCallback: successCallback,

		errChan: errChan,
		posChan: positionChan,
	}

	return w
}
