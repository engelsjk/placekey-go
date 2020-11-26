package pkapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

const (
	libraryVersion = "0.0.1"
	userAgent      = "placekey-go/" + libraryVersion

	defaultBaseURL = "https://api.placekey.io/"
	mediaType      = "application/json"

	headerRateLimitSecond     = "X-RateLimit-Limit-second"
	headerRateRemainingSecond = "X-RateLimit-Remaining-second"
	headerRateLimitMinute     = "X-RateLimit-Limit-minute"
	headerRateRemainingMinute = "X-RateLimit-Remaining-minute"

	defaultMaxRetries = 20
)

type Client struct {
	client    *http.Client
	apiKey    string
	BaseURL   *url.URL
	UserAgent string
	Rate      Rate
	ratemtx   sync.Mutex

	// Services
	SingleLocation SingleLocationService
	Bulk           BulkService

	onRequestCompleted RequestCompletionCallback
}

type RequestCompletionCallback func(*http.Request, *http.Response)

type Response struct {
	*http.Response
	Rate
}

type ErrorResponse struct {
	Response  *http.Response
	Message   string `json:"message"`
	RequestID string `jsonL"request_id`
}

type Rate struct {
	LimitSec     int `json:"limit_second"`
	RemainingSec int `json:"remaining_second"`
	LimitMin     int `json:"limit_minute"`
	RemainingMin int `json:"remaining_minute"`
}

func NewClient(apiKey string) *Client {
	httpClient := http.DefaultClient

	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{client: httpClient, apiKey: apiKey, BaseURL: baseURL, UserAgent: userAgent}
	c.SingleLocation = &SingleLocationServiceOp{client: c}
	c.Bulk = &BulkServiceOp{client: c}

	return c
}

// Options

func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("apikey", c.apiKey)
	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

func (c *Client) OnRequestCompleted(rc RequestCompletionCallback) {
	c.onRequestCompleted = rc
}

func (c *Client) GetRate() Rate {
	c.ratemtx.Lock()
	defer c.ratemtx.Unlock()
	return c.Rate
}

func newResponse(r *http.Response) *Response {
	response := Response{Response: r}
	response.populateRate()

	return &response
}

func (r *Response) populateRate() {
	if limitSec := r.Header.Get(headerRateLimitSecond); limitSec != "" {
		r.Rate.LimitSec, _ = strconv.Atoi(limitSec)
	}
	if remainingSec := r.Header.Get(headerRateRemainingSecond); remainingSec != "" {
		r.Rate.RemainingSec, _ = strconv.Atoi(remainingSec)
	}
	if limitMin := r.Header.Get(headerRateLimitMinute); limitMin != "" {
		r.Rate.LimitMin, _ = strconv.Atoi(limitMin)
	}
	if remainingMin := r.Header.Get(headerRateRemainingMinute); remainingMin != "" {
		r.Rate.RemainingMin, _ = strconv.Atoi(remainingMin)
	}
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := DoRequestWithClient(ctx, c.client, req)
	if err != nil {
		return nil, err
	}
	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp)
	}

	defer func() {
		const maxBodySlurpSize = 2 << 10
		if resp.ContentLength == -1 || resp.ContentLength <= maxBodySlurpSize {
			io.CopyN(ioutil.Discard, resp.Body, maxBodySlurpSize)
		}

		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	response := newResponse(resp)
	c.ratemtx.Lock()
	c.Rate = response.Rate
	c.ratemtx.Unlock()

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return nil, err
			}
		}
	}

	return response, err
}

func DoRequestWithClient(
	ctx context.Context,
	client *http.Client,
	req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return client.Do(req)
}

func (r *ErrorResponse) Error() string {
	if r.RequestID != "" {
		return fmt.Sprintf("%v %v: %d (request %q) %v",
			r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.RequestID, r.Message)
	}
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message)
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Message = string(data)
		}
	}

	return errorResponse
}
