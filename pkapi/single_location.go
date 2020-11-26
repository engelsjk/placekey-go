package pkapi

import (
	"context"
	"net/http"
)

const (
	singleLocationPath = "v1/placekey"
)

type SingleLocationService interface {
	// ... list methods
}

type SingleLocationServiceOp struct {
	client *Client
}

var _ SingleLocationService = &SingleLocationServiceOp{}

type SingleLocation struct {
	QueryID  string `json:"query_id"`
	Placekey string `json:"placekey"`
}

type SingleLocationRequest struct {
	Query   Query   `json:"query"`
	Options Options `json:"options"`
}

type Query struct {
	QueryID        string  `json:"query_id"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	LocationName   string  `json:"location_name"`
	StreetAddress  string  `json:"street_address"`
	City           string  `json:"city"`
	Region         string  `json:"region"`
	PostalCode     string  `json:"postal_code"`
	ISOCountryCode string  `json:"iso_country_code"`
}

type Options struct {
	StrictAddressMatch bool `json:"strict_address_match"`
}

func (svc *SingleLocationServiceOp) Get(ctx context.Context, query Query) (*SingleLocation, *Response, error) {
	req, err := svc.client.NewRequest(ctx, http.MethodPost, singleLocationPath, query)
	if err != nil {
		return nil, nil, err
	}

	sl := new(SingleLocation)

	resp, err := svc.client.Do(ctx, req, sl)
	if err != nil {
		return nil, resp, err
	}

	return sl, resp, nil
}
