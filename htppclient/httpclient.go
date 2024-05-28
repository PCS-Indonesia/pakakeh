package httpclient

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type CustomHttpClient struct {
	client  DoReq
	timeout time.Duration
	// for retry mechanism
	retrier    Retriable
	retryCount int
}

const (
	defaultRetryCount = 0
	defaultTimeout    = 30 * time.Second
)

var _ HttpMethods = (*CustomHttpClient)(nil)

// NewClient returns a new instance of http Client
func NewClient(opts ...Option) *CustomHttpClient {
	client := CustomHttpClient{
		timeout:    defaultTimeout,
		retryCount: defaultRetryCount,
		retrier:    NewNoRetrier(),
	}

	for _, opt := range opts {
		opt(&client)
	}

	if client.client == nil {
		client.client = &http.Client{
			Timeout: client.timeout,
		}
	}

	return &client
}

// Get makes a HTTP GET request to provided URL
func (c *CustomHttpClient) Get(url string, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return response, errors.Wrap(err, "GET - request process failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Post makes a HTTP POST request to provided URL and requestBody
func (c *CustomHttpClient) Post(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return response, errors.Wrap(err, "POST - request process failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Put makes a HTTP PUT request to provided URL and requestBody
func (c *CustomHttpClient) Put(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return response, errors.Wrap(err, "PUT - request process failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Patch makes a HTTP PATCH request to provided URL and requestBody
func (c *CustomHttpClient) Patch(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return response, errors.Wrap(err, "PATCH - request process failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Delete makes a HTTP DELETE request with provided URL
func (c *CustomHttpClient) Delete(url string, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return response, errors.Wrap(err, "DELETE - request process failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Do makes an HTTP request with `http.Do`
func (c *CustomHttpClient) Do(request *http.Request) (*http.Response, error) {
	var bodyReader *bytes.Reader

	if request.Body != nil {
		reqData, err := io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}

		bodyReader = bytes.NewReader(reqData)
		request.Body = io.NopCloser(bodyReader) // prevents closing the body between retries
	}

	multiErr := &Errors{}
	var response *http.Response

	for i := 0; i <= c.retryCount; i++ {
		if response != nil {
			response.Body.Close()
		}

		var err error
		response, err = c.client.Do(request)
		if bodyReader != nil {
			// Reset the body reader after the request since at this point it's already read
			// Note that it's safe to ignore the error here since the 0,0 position is always valid
			_, _ = bodyReader.Seek(0, 0)
		}

		if err != nil {
			multiErr.Push(err.Error())

			backoffTime := c.retrier.NextInterval(i)
			time.Sleep(backoffTime)

			continue
		}

		if response.StatusCode >= http.StatusInternalServerError {
			backoffTime := c.retrier.NextInterval(i)
			time.Sleep(backoffTime)
			continue
		}

		multiErr = &Errors{} // Clear ALL errors if any iteration process succeeds
		break
	}

	return response, multiErr.HasError()
}
