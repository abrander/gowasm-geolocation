# gowasm-geolocation

Geolocation is an idiomatic and intuitive Go package for using
the browser's geolocation API.

[![GoDoc][1]][2]

[1]: https://godoc.org/github.com/abrander/gowasm-geolocation?status.svg
[2]: https://godoc.org/github.com/abrander/gowasm-geolocation

## Examples

### Getting the device positon:

```go
package main

import (
	"fmt"

	"github.com/abrander/gowe/geolocation"
)

func main() {
	options := &geolocation.PositionOptions{
		HighAccuracy: true,
	}

	pos, err := geolocation.CurrentPosition(options)
	if err != nil {
		fmt.Printf("Ohh no: %s\n", err.Error())
	}

	fmt.Printf("Got position at %s: %+v\n", pos.Timestamp.String(), pos.Coords)
}
```


### Watching the device postion:

```go
package main

import (
	"fmt"

	"github.com/abrander/gowe/geolocation"
)

func main() {
	w := geolocation.WatchPosition(nil)

	for {
		pos, err := w.Next()
		if err != nil {
			if err.Temporary() {
				continue
			}

			fmt.Printf("Something went wrong: %s\n", err.Error())

			break
		}

		fmt.Printf("Got new position at %s: %+v\n", pos.Timestamp.String(), pos.Coords)
	}
	w.Close()
}
```

### Integrating into a select style main loop

```go
package main

import (
	"fmt"

	"github.com/abrander/gowe/geolocation"
)

func main() {
	options := &geolocation.PositionOptions{
		MaximumAge: 10 * time.Second,
	}

	w := geolocation.WatchPosition(options)
	positions, locationErrors := w.Chans()

MAIN:
	for {
		select {
		case pos := <-positions:
			fmt.Printf("Got position at %s: %+v\n", pos.Timestamp.String(), pos.Coords)
		case err := <-locationErrors:
			fmt.Printf("Ohh no: %s\n", err.Error())
			if !err.Temporary() {
				break MAIN
			}
		}
	}

	w.Close()
}
```
