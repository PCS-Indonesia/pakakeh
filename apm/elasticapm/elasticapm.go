package elasticapm

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmhttp"
	"go.elastic.co/apm/stacktrace"
)

func init() {
	stacktrace.RegisterLibraryPackage(
		"github.com/gin-gonic",
		"github.com/gin-contrib",
	)
}

// struct apmMiddleware
type apmMiddleware struct {
	engine         *gin.Engine
	tracer         *apm.Tracer
	requestIgnorer apmhttp.RequestIgnorerFunc
}

// ApmMiddleware returns a new Gin middleware handler for tracing (specific for Gin)
// requests, run time process and reporting errors. This package is modification
// from go elastic apm package for internal PCS Payment use
func ApmMiddleware(engine *gin.Engine, opts ...Option) gin.HandlerFunc {
	am := &apmMiddleware{
		engine: engine,
		tracer: apm.DefaultTracer, // set as default tracer if caller does not define with WithTracer option
	}

	for _, opt := range opts {
		opt(am)
	}

	if am.requestIgnorer == nil {
		// NewDynamicServerRequestIgnorer returns the RequestIgnorer to use in handler
		am.requestIgnorer = apmhttp.NewDynamicServerRequestIgnorer(am.tracer)
	}

	return am.handle
}

// getRequestName returns the transaction name for the server request
func getRequestName(c *gin.Context) string {
	if fullPath := c.FullPath(); fullPath != "" {
		return c.Request.Method + " " + fullPath
	}
	return apmhttp.ServerRequestName(c.Request)
}

// will handle middleware when called
func (am *apmMiddleware) handle(c *gin.Context) {
	if !am.tracer.Recording() || am.requestIgnorer(c.Request) {
		c.Next()

		return
	}

	requestName := getRequestName(c)

	// set incoming transaction with Body you've defined
	tx, body, req := apmhttp.StartTransactionWithBody(am.tracer, requestName, c.Request)
	defer tx.End()

	c.Request = req

	defer func() {
		if v := recover(); v != nil {
			if !c.Writer.Written() {
				c.AbortWithStatus(http.StatusInternalServerError)
			} else {
				c.Abort()
			}

			e := am.tracer.Recovered(v)
			e.SetTransaction(tx)
			setContext(&e.Context, c, body)
			e.Send()
		}

		c.Writer.WriteHeaderNow()
		tx.Result = apmhttp.StatusCodeResult(c.Writer.Status())

		if tx.Sampled() {
			setContext(&tx.Context, c, body)
		}

		for _, err := range c.Errors {
			e := am.tracer.NewError(err.Err)
			e.SetTransaction(tx)

			setContext(&e.Context, c, body)

			e.Handled = true
			e.Send()
		}

		body.Discard()
	}()

	// call next handler
	c.Next()
}

// setContext set all apmContext value (PCS use Gin as framework, etc)
func setContext(ctx *apm.Context, c *gin.Context, body *apm.BodyCapturer) {
	ctx.SetFramework("gin", gin.Version) // we are using Gin by default
	ctx.SetHTTPRequest(c.Request)
	ctx.SetHTTPRequestBody(body)
	ctx.SetHTTPResponseHeaders(c.Writer.Header())
	ctx.SetHTTPStatusCode(c.Writer.Status())
}
