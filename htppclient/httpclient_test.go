package httpclient

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Ini buat mock response body dalam string
func mockRespBody(t *testing.T, response *http.Response) string {
	if response.Body != nil {
		defer response.Body.Close()
	}

	mockRespBody, err := io.ReadAll(response.Body)
	require.NoError(t, err, "success to read response body")

	return string(mockRespBody)
}

// Unit Test defined below
func TestHTTPClientDoSuccess(t *testing.T) {
	client := NewClient(WithTimeout(10 * time.Millisecond))

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok gas" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(testHandler))
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "en")

	response, err := client.Do(req)
	require.NoError(t, err, "success to make a GET request")

	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, "{ \"response\": \"ok gas\" }", string(body))
}

func TestHTTPClientGetSuccess(t *testing.T) {
	client := NewClient(WithTimeout(10 * time.Millisecond))

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok gas" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")

	response, err := client.Get(server.URL, headers)
	require.NoError(t, err, "success to make a GET request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok gas\" }", mockRespBody(t, response))
}

func TestHTTPClientDeleteSuccess(t *testing.T) {
	client := NewClient(WithTimeout(10 * time.Millisecond))

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok gas" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")

	response, err := client.Delete(server.URL, headers)
	require.NoError(t, err, "success to make a DELETE request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok gas\" }", mockRespBody(t, response))
}

func TestHTTPClientPostSuccess(t *testing.T) {
	client := NewClient(WithTimeout(10 * time.Millisecond))

	requestBodyString := `{ "name": "Mbah Surip" }`

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		rBody, err := io.ReadAll(r.Body)
		require.NoError(t, err, "success to extract request body")

		assert.Equal(t, requestBodyString, string(rBody))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	requestBody := bytes.NewReader([]byte(requestBodyString))

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")

	response, err := client.Post(server.URL, requestBody, headers)
	require.NoError(t, err, "success to make a POST request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", mockRespBody(t, response))
}

func TestHTTPClientPutSuccess(t *testing.T) {
	client := NewClient(WithTimeout(10 * time.Millisecond))

	requestBodyString := `{ "name": "Alfarizi Pusing" }`

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		rBody, err := io.ReadAll(r.Body)
		require.NoError(t, err, "success to extract request body")

		assert.Equal(t, requestBodyString, string(rBody))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	requestBody := bytes.NewReader([]byte(requestBodyString))

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")

	response, err := client.Put(server.URL, requestBody, headers)
	require.NoError(t, err, "success to make a PUT request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", mockRespBody(t, response))
}

func TestHTTPClientPatchSuccess(t *testing.T) {
	client := NewClient(WithTimeout(10 * time.Millisecond))

	requestBodyString := `{ "name": "Test Pak El" }`

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		rBody, err := io.ReadAll(r.Body)
		require.NoError(t, err, "success extract request body")

		assert.Equal(t, requestBodyString, string(rBody))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	requestBody := bytes.NewReader([]byte(requestBodyString))

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")

	response, err := client.Patch(server.URL, requestBody, headers)
	require.NoError(t, err, "success to make a PATCH request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", mockRespBody(t, response))
}

func TestHTTPClientGetRetriesOnFailure(t *testing.T) {
	count := 0
	numOfRetries := 3
	numOfCalls := numOfRetries + 1
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	// Coba retry
	client := NewClient(
		WithTimeout(10*time.Millisecond),
		WithRetryCount(numOfRetries),
		WithRetrier(NewRetrier(NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "response": "something's wrong" }`))
		count++
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	response, err := client.Get(server.URL, http.Header{})
	require.NoError(t, err, "failed to make GET request")

	require.Equal(t, http.StatusInternalServerError, response.StatusCode)
	require.Equal(t, "{ \"response\": \"something's wrong\" }", mockRespBody(t, response))

	assert.Equal(t, numOfCalls, count)
}

func TestHTTPClientPostRetriesOnFailure(t *testing.T) {
	count := 0
	numOfRetries := 3
	numOfCalls := numOfRetries + 1
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithTimeout(10*time.Millisecond),
		WithRetryCount(numOfRetries),
		WithRetrier(NewRetrier(NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "response": "error make GET request" }`))
		count++ // set count until maximum
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	response, err := client.Post(server.URL, strings.NewReader("a=1"), http.Header{})
	require.NoError(t, err, "Failed to make GET request after retry")

	require.Equal(t, http.StatusInternalServerError, response.StatusCode)
	require.Equal(t, "{ \"response\": \"Failed to make GET request\" }", mockRespBody(t, response))

	assert.Equal(t, numOfCalls, count)
}

func TestHTTPClientGetReturnsNoErrorsIfRetriesFailWith5xx(t *testing.T) {
	count := 0
	numOfRetries := 2
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithTimeout(10*time.Millisecond),
		WithRetryCount(numOfRetries),
		WithRetrier(NewRetrier(NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "response": "something's wrong" }`))
		count++
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	response, err := client.Get(server.URL, http.Header{})
	require.NoError(t, err)

	require.Equal(t, http.StatusInternalServerError, response.StatusCode)
	require.Equal(t, numOfRetries+1, count)
	require.Equal(t, "{ \"response\": \"something's wrong\" }", mockRespBody(t, response))
}

func TestHTTPClientGetReturnsNoErrorsIfRetrySucceeds(t *testing.T) {
	count := 0
	countWhenCallSucceeds := 2
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithTimeout(10*time.Millisecond),
		WithRetryCount(3),
		WithRetrier(NewRetrier(NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		if count == countWhenCallSucceeds {
			w.Write([]byte(`{ "response": "success" }`))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{ "response": "something's wrong" }`))
		}
		count++
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	response, err := client.Get(server.URL, http.Header{})
	require.NoError(t, err, "success to make GET request")

	require.Equal(t, countWhenCallSucceeds+1, count)
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.Equal(t, "{ \"response\": \"success\" }", mockRespBody(t, response))
}

func TestHTTPClientGetReturnsErrorOnClientCallFailure(t *testing.T) {
	client := NewClient(WithTimeout(10 * time.Millisecond))

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	server.URL = "" // Mock Invalid URL
	defer server.Close()

	response, err := client.Get(server.URL, http.Header{})
	require.Error(t, err, "error make GET request")

	require.Nil(t, response)

	assert.Contains(t, err.Error(), "unsupported protocol scheme")
}

func TestHTTPClientGetReturnsNoErrorOn5xxFailure(t *testing.T) {
	client := NewClient(WithTimeout(10 * time.Millisecond))

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "response": "something went wrong" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	response, err := client.Get(server.URL, http.Header{})
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, response.StatusCode)

}

func TestHTTPClientGetReturnsErrorOnFailure(t *testing.T) {
	client := NewClient(WithTimeout(10 * time.Millisecond))

	response, err := client.Get("invalid_url", http.Header{})
	assert.Contains(t, err.Error(), "unsupported protocol scheme")
	assert.Nil(t, response)
}

type myCustomHTTPClient struct {
	client http.Client
}

func (c *myCustomHTTPClient) Do(request *http.Request) (*http.Response, error) {
	request.Header.Set("test", "tist")
	return c.client.Do(request)
}

func TestCustomHTTPClientHeaderSuccess(t *testing.T) {
	client := NewClient(
		WithTimeout(10*time.Millisecond),
		WithHTTPClient(&myCustomHTTPClient{
			client: http.Client{Timeout: 25 * time.Millisecond}}),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "tist", r.Header.Get("test"))
		assert.NotEqual(t, "toss", r.Header.Get("tist"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	response, _ := client.Do(req)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, "{ \"response\": \"ok\" }", string(body))
}
