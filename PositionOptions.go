package geolocation

import (
	"syscall/js"
	"time"
)

type PositionOptions struct {
	// HighAccuracy provides a hint that the application would like
	// to receive the best possible results. This may result in
	// slower response times or increased power consumption. The
	// user might also deny this capability, or the device might
	// not be able to provide more accurate results than if the
	// flag wasn't specified.
	HighAccuracy bool

	// Timeout denotes the maximum length of time that is allowed
	// to pass until a position or an error is returned.
	Timeout time.Duration

	// MaximumAge indicates that the application is willing to
	// accept a cached position whose age is no greater than the
	// specified time.
	MaximumAge time.Duration
}

func (o *PositionOptions) jsValue() js.Value {
	proxy := make(map[string]interface{})

	if o != nil {
		if o.HighAccuracy {
			proxy["enableHighAccuracy"] = true
		}

		if o.Timeout > 0 {
			proxy["timeout"] = o.Timeout.Milliseconds()
		}

		proxy["maximumAge"] = o.MaximumAge.Milliseconds()
	}

	return js.ValueOf(proxy)
}
