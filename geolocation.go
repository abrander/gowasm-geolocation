// Package geolocation provides a convenient idiomatic wrapper
// for the browser geolocation API.
package geolocation

import (
	"syscall/js"
)

func geolocation() js.Value {
	return js.Global().Get("navigator").Get("geolocation")
}
