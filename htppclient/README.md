# Pakakeh httpclient

## About
<p>Pakakeh httpclient is an custom HTTP client that helps your application make a large number of requests, at scale i hope :D  </p>

It is built using native go library `net/http`. With this you can :
- Create clients with different timeouts for every request or you can use default configuration (see example in unit test `httpclient_test.go` file)
- Set synchronous retries for each request, with the option of setting your own retry strategies.

## How to use

This package can be used by adding the following import `httpclient` statement to your `.go` files.

```go
import "github.com/PCS-Indonesia/pakakeh/httpclient" 
```

## How to create custom httpclient
To create custom httpclient using this package, you need to declare the client by using `NewClient(opt options...)` function with variadic options param such as :
- WithTimeout
- WithRetrier
- WithHTTPClient
- WithRetryCount
<br></br>

#### example making simple of GET Request
```go
// Create a new HTTP client with a default timeout
timeout := 10 * time.Second
client := httpclient.NewClient(httpclient.WithHTTPTimeout(timeout))

// Use the clients GET method to create and execute the request
res, err := client.Get("http://google.com", nil)
if err != nil{
	panic(err)
}

// pakakeh httpclient returns the native *http.Response 
respBody, err := io.ReadAll(res.Body)
fmt.Println(string(respBody))
```
<br></br>
#### or you can leverage Go built-in `*http.NewRequest` with defined `http.Do` interface :

```go
// Create an http.Request instance
req, _ := http.NewRequest(http.MethodGet, "http://google.com", nil)

res, err := client.Do(req)
if err != nil {
	panic(err)
}

respBody, err := io.ReadAll(res.Body)
fmt.Println(string(respBody))

```
</br>

#### Create HTTP client with a retry mechanism
If you are familiar with jitter or other retry mechanism in http client, then this will be easy to understand about interval coefficients. Also, if you implementing this in API, you need to set Env to store the time interval. For simplicity, i will show example using hardcoded values.

```go
// backoffInterval increases the backoff at a constant rate. set this first before it's too late
backoffInterval := 2 * time.Second
// Define a maximum jitter interval.
jitterInterval := 1 * time.Second

backoff := httpclient.NewConstantBackoff(backoffInterval, jitterInterval)

// retry mechanism with the backoff
retryMech := httpclient.NewRetrier(backoff)

timeout := 30 * time.Second

// sets new client, the retry mechanism, and the number of times you like to retry
client := httpclient.NewClient(
	httpclient.WithTimeout(timeout), // Timeout for each request
	httpclient.WithRetrier(retryMech), // retry mechanisms method
	httpclient.WithRetryCount(3), // retry count
)

// The rest is same as the first above example
..
...
....
```

</br>

If you are familiar with jitter or other retry mechanism in http client, then this will be easy to understand about interval coefficients. Also, if you implementing this in API, you need to set Env to store the time interval. For simplicity, i will show example using hardcoded values.

```go
initTimeout := 2 * time.Second // initial TO
maxTimeout := 4 * time.Second // max TO
exponentFactor := 2 // exponent factor
jitterInterval := 1 * time.Second // Define a maximum jitter interval. 

backoff := httpclient.NewExponentialBackoff(initTimeout, maxTimeout, exponentFactor, jitterInterval)

// retry mechanism with the backoff
retryMech := httpclient.NewRetrier(backoff)

timeout := 30 * time.Second

// sets new client, the retry mechanism, and the number of times you like to retry
client := httpclient.NewClient(
	httpclient.WithTimeout(timeout), // Timeout for each request
	httpclient.WithRetrier(retryMech), // retry mechanisms method
	httpclient.WithRetryCount(3), // retry count
)

// The rest is same as the first above example
..
...
....
```

</br>
Not only that, Pakakeh httpclient also supports custom retry strategies (if you want). To implement, you must implement `Backoff` interface

```go
type Backoff interface {
	Next(retry int) time.Duration
}

type linearBackoff struct {
	backoffInterval int
}

func (lb *linearBackoff) Next(retry int) time.Duration{
	if retry <= 0 {
		return 0 * time.Millisecond
	}
	return time.Duration(retry * lb.backoffInterval) * time.Millisecond
}

backoff := &linearBackoff{1000} // in millisecond
retryMech := httpclient.NewRetrier(backoff)

timeout := 10 * time.Second
// Create a new client, sets the retry mechanism, and the number of times you would like to retry
client := httpclient.NewClient(
	httpclient.WithHTTPTimeout(timeout),
	httpclient.WithRetrier(retryMech),
	httpclient.WithRetryCount(4),
)

// The rest is same as the first above example
..
...
....
```