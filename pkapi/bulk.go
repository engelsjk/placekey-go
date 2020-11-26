package pkapi

import (
	"context"
	"net/http"
)

const (
	bulkPath = "v1/placekeys"
)

type BulkService interface {
	Get(context.Context, *BulkRequest) (*Bulk, *Response, error)
}

type BulkServiceOp struct {
	client *Client
}

var _ BulkService = &BulkServiceOp{}

type Bulk []SingleLocation

type BulkRequest struct {
	Queries []Query `json:"queries"`
}

func (svc *BulkServiceOp) Get(ctx context.Context, request *BulkRequest) (*Bulk, *Response, error) {
	req, err := svc.client.NewRequest(ctx, http.MethodPost, singleLocationPath, request)
	if err != nil {
		return nil, nil, err
	}

	b := new(Bulk)

	resp, err := svc.client.Do(ctx, req, b)
	if err != nil {
		return nil, resp, err
	}

	return b, resp, nil
}
