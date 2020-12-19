package geolocation

import (
	"math"
	"syscall/js"
)

// Coords represents a place on Earth. The geographic coordinate
// reference system used is the World Geodetic System (2d) [WGS84].
type Coords struct {
	// Latitude and Longitude are geographic coordinates
	// specified in decimal degrees.
	Latitude  float64
	Longitude float64

	// Accuracy denotes the accuracy level of the latitude and
	// longitude coordinates in meters.
	Accuracy float64

	altitude         *float64
	altitudeAccuracy *float64
	heading          *float64
	speed            *float64
}

const (
	// ErrUnknown will be returned if the value is not supplied by
	// the implementation.
	ErrUnknown = Error("value unknown")

	// ErrStationary will be returned if the device velocity is
	// requested when the device is stationary.
	ErrStationary = Error("device stationary")
)

func toCoord(value js.Value) *Coords {
	nilOrFloat := func(val js.Value) *float64 {
		var f float64

		switch {
		case val.IsNull():
			return nil

		case val.IsNaN():
			f = math.NaN()

		default:
			f = val.Float()
		}

		return &f
	}

	coords := &Coords{
		Latitude:         value.Get("latitude").Float(),
		Longitude:        value.Get("longitude").Float(),
		Accuracy:         value.Get("accuracy").Float(),
		altitude:         nilOrFloat(value.Get("altitude")),
		altitudeAccuracy: nilOrFloat(value.Get("altitudeAccuracy")),
		heading:          nilOrFloat(value.Get("heading")),
		speed:            nilOrFloat(value.Get("speed")),
	}

	return coords
}

// Altitude denotes the height of the position, specified
// in meters above the [WGS84] ellipsoid. Altitude may not be
// avilable on all platforms at all times. ErrUnknown will
// be returned, if the altitude is not known or unsupported
// by the host device.
func (c *Coords) Altitude() (float64, error) {
	if c.altitude != nil {
		return *c.altitude, nil
	}

	return math.NaN(), ErrUnknown
}

// AltitudeAccuracy is specified in meters. This may return an
// ErrUnknown error if the accuracy is unavailable og unsupported
// on the host device.
func (c *Coords) AltitudeAccuracy() (float64, error) {
	if c.altitudeAccuracy != nil {
		return *c.altitudeAccuracy, nil
	}

	return math.NaN(), ErrUnknown
}

// Heading is the direction of travel of the device
// and is specified in degrees, where 0° ≤ Heading < 360°,
// counting clockwise relative to the true north.
// If this is unknown or unsupported, it will return
// an ErrUnknown error.
// If the host device is not moving, it will return
// ErrStationary.
func (c *Coords) Heading() (float64, error) {
	if c.heading != nil {
		if math.IsNaN(*c.heading) {
			return math.NaN(), ErrStationary
		}

		return *c.heading, nil
	}

	return math.NaN(), ErrUnknown
}

// Speed is the magnitude of the horizontal component
// of the hosting device's current velocity and is specified
// in meters per second.
// If the velocity is not known or unsupported, an ErrUnknown
// error will be returned.
func (c *Coords) Speed() (float64, error) {
	if c.speed != nil {
		return *c.speed, nil
	}

	return math.NaN(), ErrUnknown
}
