# placekey-go

A Go library for working with [Placekeys](https://www.placekey.io/). This is an **unofficial** port of the Python library [placekey-py](https://github.com/Placekey/placekey-py).

## Status

This library is mostly complete. Some geospatial features are currently under development. An interface to the [Placekeys API](https://docs.placekey.io) may also be added at some point.

## Installation

```
go get -u github.com/engelsjk/placekey-go
```

## Usage

```
> placekey.FromGeo(37.23712, -115.80187)
@5ys-rsx-4jv

> placekey.ToH3("@5yv-j8h-3nq")
8a2986b843b7fff

> placekey.Distance("@5ys-rsx-4jv", "@5yv-j8h-3nq")
138681.552855
```

### Dependencies

* [uber/h3-go](github.com/uber/h3-go)
