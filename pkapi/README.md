# placekey-go/pkapi

An API client library for the [Placekey API](https://docs.placekey.io/).

## Usage

```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/engelsjk/placekey-go/pkapi"
)

func main() {

  api := pkapi.NewClient(os.Getenv("PLACEKEY_API_KEY"))

  ctx := context.Background()
  ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
  defer cancel()

  req := &pkapi.SingleLocationRequest{
    Query: pkapi.Query{
      Latitude:  37.23712,
      Longitude: -115.80187,
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

  fmt.Printf("rate: %+v\n", api.GetRate())
  // rate: {LimitSec:100 RemainingSec:99 LimitMin:1000 RemainingMin:999}

  fmt.Printf("response: %s\n", string(b))
  // response: {"query_id":"0","placekey":"@5vg-82n-kzz"}
}

```
