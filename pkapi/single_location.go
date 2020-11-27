package pkapi

import (
	"context"
	"net/http"
)

const (
	singleLocationPath = "v1/placekey"
)

type SingleLocationService interface {
	Get(context.Context, *SingleLocationRequest) (*SingleLocation, *Response, error)
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
	Query   Query    `json:"query"`
	Options *Options `json:"options,omitempty"`
}

type Query struct {
	QueryID        string  `json:"query_id,omitempty"`
	Latitude       float64 `json:"latitude,omitempty"`
	Longitude      float64 `json:"longitude,omitempty"`
	LocationName   string  `json:"location_name,omitempty"`
	StreetAddress  string  `json:"street_address,omitempty"`
	City           string  `json:"city,omitempty"`
	Region         string  `json:"region,omitempty"`
	PostalCode     string  `json:"postal_code,omitempty"`
	ISOCountryCode string  `json:"iso_country_code,omitempty"`
}

type Options struct {
	StrictAddressMatch bool `json:"strict_address_match,omitempty"`
	StrictNameMatch    bool `json:"strict_name_match,omitempty"`
}

// Get sends a Singe Location request to the Placekey API and returns a Placekey responses.
func (svc *SingleLocationServiceOp) Get(ctx context.Context, request *SingleLocationRequest) (*SingleLocation, *Response, error) {
	req, err := svc.client.NewRequest(ctx, http.MethodPost, singleLocationPath, request)
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
