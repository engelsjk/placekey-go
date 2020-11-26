package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/engelsjk/placekey-go/pkapi"
)

func main() {

	api := pkapi.NewClient(os.Getenv("PLACEKEY_API_KEY"))

	ExampleSingleCoordinates(api)
	ExampleSingleCoordinatesWithQueryID(api)
	ExampleSingleAddress(api)
	ExampleSingleAddressStrictAddress(api)
	ExampleSinglePOI(api)
	ExampleSinglePOIStrictName(api)
	ExampleBulk(api)
}

func ExampleSingleCoordinates(api *pkapi.Client) {

	fmt.Println("*** single_location_coordinates ***")

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &pkapi.SingleLocationRequest{
		Query: pkapi.Query{
			Latitude:  37.7371,
			Longitude: -122.44283,
		},
	}

	sl, resp, err := api.SingleLocation.Get(ctx, req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	b, err := json.Marshal(sl)
	if err != nil {
		panic(err)
	}

	fmt.Printf("...rate: %+v\n", api.GetRate())
	fmt.Printf("...response: %s\n", string(b))
}

func ExampleSingleCoordinatesWithQueryID(api *pkapi.Client) {

	fmt.Println("*** single_location_coordinates_with_query ***")

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &pkapi.SingleLocationRequest{
		Query: pkapi.Query{
			QueryID:   "thisiscustom",
			Latitude:  37.7371,
			Longitude: -122.44283,
		},
	}

	sl, resp, err := api.SingleLocation.Get(ctx, req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	b, err := json.Marshal(sl)
	if err != nil {
		panic(err)
	}

	fmt.Printf("...rate: %+v\n", api.GetRate())
	fmt.Printf("...response: %s\n", string(b))
}

func ExampleSingleAddress(api *pkapi.Client) {

	fmt.Println("*** single_location_address ***")

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &pkapi.SingleLocationRequest{
		Query: pkapi.Query{
			StreetAddress:  "1 Dr Carlton B Goodlett Pl",
			City:           "San Francisco",
			Region:         "CA",
			PostalCode:     "94102",
			ISOCountryCode: "US",
		},
	}

	sl, resp, err := api.SingleLocation.Get(ctx, req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	b, err := json.Marshal(sl)
	if err != nil {
		panic(err)
	}

	fmt.Printf("...rate: %+v\n", api.GetRate())
	fmt.Printf("...response: %s\n", string(b))
}

func ExampleSingleAddressStrictAddress(api *pkapi.Client) {

	fmt.Println("*** single_location_address_strict ***")

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &pkapi.SingleLocationRequest{
		Query: pkapi.Query{
			StreetAddress:  "598 Portola Dr",
			City:           "San Francisco",
			Region:         "CA",
			PostalCode:     "94131",
			ISOCountryCode: "US",
		},
		Options: &pkapi.Options{
			StrictAddressMatch: true,
		},
	}

	sl, resp, err := api.SingleLocation.Get(ctx, req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	b, err := json.Marshal(sl)
	if err != nil {
		panic(err)
	}

	fmt.Printf("...rate: %+v\n", api.GetRate())
	fmt.Printf("...response: %s\n", string(b))
}

func ExampleSinglePOI(api *pkapi.Client) {

	fmt.Println("*** single_location_poi ***")

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &pkapi.SingleLocationRequest{
		Query: pkapi.Query{
			LocationName:   "San Francisco City Hall",
			StreetAddress:  "1 Dr Carlton B Goodlett Pl",
			City:           "San Francisco",
			Region:         "CA",
			PostalCode:     "94102",
			ISOCountryCode: "US",
		},
	}

	sl, resp, err := api.SingleLocation.Get(ctx, req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	b, err := json.Marshal(sl)
	if err != nil {
		panic(err)
	}

	fmt.Printf("...rate: %+v\n", api.GetRate())
	fmt.Printf("...response: %s\n", string(b))
}

func ExampleSinglePOIStrictName(api *pkapi.Client) {

	fmt.Println("*** single_location_poi_strict ***")

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &pkapi.SingleLocationRequest{
		Query: pkapi.Query{
			LocationName:   "San Francisco City Hall",
			StreetAddress:  "1 Dr Carlton B Goodlett Pl",
			City:           "San Francisco",
			Region:         "CA",
			PostalCode:     "94102",
			ISOCountryCode: "US",
		},
		Options: &pkapi.Options{
			StrictNameMatch: true,
		},
	}

	sl, resp, err := api.SingleLocation.Get(ctx, req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	b, err := json.Marshal(sl)
	if err != nil {
		panic(err)
	}

	fmt.Printf("...rate: %+v\n", api.GetRate())
	fmt.Printf("...response: %s\n", string(b))
}

func ExampleBulk(api *pkapi.Client) {

	fmt.Println("*** bulk ***")

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &pkapi.BulkRequest{Queries: []pkapi.Query{
		{
			StreetAddress:  "1543 Mission Street, Floor 3",
			City:           "San Francisco",
			Region:         "CA",
			PostalCode:     "94105",
			ISOCountryCode: "US",
		},
		{
			QueryID:        "thisqueryidaloneiscustom",
			LocationName:   "Twin Peaks Petroleum",
			StreetAddress:  "598 Portola Dr",
			City:           "San Francisco",
			Region:         "CA",
			PostalCode:     "94131",
			ISOCountryCode: "US",
		},
		{
			Latitude:  37.7371,
			Longitude: -122.44283,
		},
	}}

	sl, resp, err := api.Bulk.Get(ctx, req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	b, err := json.Marshal(sl)
	if err != nil {
		panic(err)
	}

	fmt.Printf("...rate: %+v\n", api.GetRate())
	fmt.Printf("...response: %s\n", string(b))
}
