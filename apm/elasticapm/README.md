## Pakakeh ElasticAPM Middleware

### About
This package provides APM Middleware for Gin-Gonic API functionality using elastic search APM. 
It is built using native external go library `go.elastic.co/apm` with some modification inside.
You can use spanner to embed this APM agent to your client and wrap the elasticsearch APM functionality 
with your specific function in your Gin API (by default PCS use GIN for api framework).


## How to use
This package can be used by adding the following import `httpclient` statement to your `.go` files.
<i>i suggest you inject this into config options for clean and better abstraction </i>

```go
import "github.com/PCS-Indonesia/pakakeh/middlewares/elasticapm" 
```

Let's dive in how to use this with some following `example API` :
1. Create an API with GET method for “/example” path. Inside the router handler, use `elasticapm` in your built-in Gin Engine `Use(someMiddleware)` method, Then we will simulate some processes in which consists of 3 functions. 
2. `processRequest` function that will sleep 15 miliseconds
3. `doSomething` function that will sleep 20 miliseconds
4. `getApaAjaFromAPI` function that will fetch example data from JSON placeholder API (https://jsonplaceholder.typicode.com/todos/1)
<br>

Since we’re using gin framework, we can utilize the `Use()` Middleware function that will automatically wrap any router handler created by us and send it to the APM server. We can continuing the span created by the middleware by golang `Context` then we create new span (extends) from it

```go
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/PCS-Indonesia/pakakeh/middlewares/elasticapm"
	"github.com/gin-gonic/gin"
	"go.elastic.co/apm"
)

func main() {
	r := gin.Default()

	tracer := elasticapm.NewTracer(
        // Better init this inside config
		"http://localhost:8200", // default url, testing in my local environment
		"",                      // set this value if APM server need ApiKey auth
		"test-API",              // API name
		"v1.0.0",                // API / Service version
	)

	r.Use(elasticapm.ApmMiddleware(r, elasticapm.WithTracer(tracer)))

	r.GET("/example", func(c *gin.Context) {
		span, ctx := apm.StartSpan(c.Request.Context(), "ExampleHandler", "request")
		defer span.End()

		processRequest(ctx)

		result, err := getApaAjaFromAPI(ctx)
		if err != nil {
			log.Println(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"result": result,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func processRequest(ctx context.Context) {
	span, ctx := apm.StartSpan(ctx, "processRequest", "custom")
	defer span.End()

	doSomething(ctx)

	// time sleep simulate some processing time
	time.Sleep(15 * time.Millisecond)
}

func doSomething(ctx context.Context) {
	span, ctx := apm.StartSpan(ctx, "doSomething", "custom")
	defer span.End()

	// time sleep simulate some processing time
	time.Sleep(20 * time.Millisecond)

}

func getApaAjaFromAPI(ctx context.Context) (map[string]interface{}, error) {
	span, ctx := apm.StartSpan(ctx, "getApaAjaFromAPI", "custom")
	defer span.End()

	var result map[string]interface{}

	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos/1")
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, err
}
```

For simplicity, we’re putting up all code in the same main.go file, but in the real case it’s better to do more clean code by dividing it into its own package and domains. As you can see, every function that you want to trace by the APM need to have golang context parameter in it’s function parameter so it can extends the parent context. Speaking of `apm.StartSpan()` method itself receives three parameters:
1. Context, commonly using the parent context in this case, or you can create a new one if it’s the first span
2. Name of the span (use anything that helps you to identify the span)
3. Type of the span (`request`, `custom`, or `anything`)

## Test APM
_I'm running the ELK APM using Docker in my local environment. This only using default configuration for the elastic stack, on real production case, more configuration will need to be tweaked (setting up elasticsearch server, nodes, indices, etc)._

Let's take a look at APM Server 
![Screenshot-2024-06-02-at-13-33-04.png](https://i.postimg.cc/3NBMgS1q/Screenshot-2024-06-02-at-13-33-04.png)

as you can see our `test-API` already shown on the Services section. If we’re having let say 20 APIs, all will be listed here if there are hits on those API, it’s automatically added when we’re adding `elasticapm` middleware.

let's move to the Transaction page to see the endpoints that hit from client
![Screenshot-2024-06-02-at-13-35-50.png](https://i.postimg.cc/ZqGhZy0w/Screenshot-2024-06-02-at-13-35-50.png)

go to below the result of the page
![Screenshot-2024-06-02-at-13-36-02.png](https://i.postimg.cc/g0QGJTtq/Screenshot-2024-06-02-at-13-36-02.png)

you can see the execution time from each function that be spawn after API is hit

On our case since the url is also being sent to the APM it’s quite simple to read (like from URL query param), but it will be useful for POST API method since the requested data is in the body