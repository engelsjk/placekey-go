# placekey-go

A Go library for working with [Placekeys](https://www.placekey.io/). This is an **unofficial** port of the Python library [placekey-py](https://github.com/Placekey/placekey-py).

## Status

This library port is mostly complete. An interface to the [Placekeys API](https://docs.placekey.io) is included. Some geospatial features are currently under development. 

## Usage

### Prerequisites

This library depends on [uber/h3-go](https://github.com/uber/h3-go) and inherits the same [prerequisites](https://github.com/uber/h3-go#prerequisites). It requires [CGO](https://golang.org/cmd/cgo/) (```CGO_ENABLED=1```) in order to be built.

> If you see errors/warnings like "build constraints exclude all Go files...", then the cgo build constraint is likely disabled; try setting CGO_ENABLED=1 environment variable for your build step.

### Installation

#### [golang/cmd/go](https://golang.org/cmd/go/)

```go
go get github.com/engelsjk/placekey-go
```

### Quickstart

```bash
import "github.com/engelsjk/placekey-go"

func ExampleFromGeo() {
    placekey.FromGeo(37.23712, -115.80187)
    // Output:
    // @5ys-rsx-4jv
}

func ExampleToH3() {
    placekey.ToH3("@5yv-j8h-3nq")
    // Output:
    // 8a2986b843b7fff
}

func ExampleDistance() {
    placekey.Distance("@5ys-rsx-4jv", "@5yv-j8h-3nq")
    // Output:
    // 138681.552855
}

func ExampleGeoJSON() {
    placekey.ToGeoJSON("@5yv-j8h-3nq")
    // Output:
    // {"type":"Feature","geometry":{"type":"Polygon","coordinates":[...]},"properties":null}
}

func ExampleDistance() {
    placekey.ToWKT("@5yv-j8h-3nq")
    // Output:
    // POLYGON(...)
}

```

### Dependencies

* [uber/h3-go](https://github.com/uber/h3-go)
