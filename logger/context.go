package logger

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type contextKey string

const (
	// RequestIDKey is the context key used to store and retrieve the request ID.
	RequestIDKey contextKey = "request_id"

	// RequestIDHeader is the HTTP header name used to pass the request ID.
	RequestIDHeader = "X-Request-Id"
)

// GetRequestID extracts the request ID from the given context.
// It supports both *gin.Context and context.Context.
// Returns an empty string if no request ID is found.
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if rID, ok := ctx.Value(RequestIDKey).(string); ok {
		return rID
	}
	return ""
}

// GetRequestIDFromGin extracts the request ID from a gin.Context.
// It looks in the request's context for the RequestIDKey.
func GetRequestIDFromGin(c *gin.Context) string {
	if c == nil {
		return ""
	}
	return GetRequestID(c.Request.Context())
}

// SetRequestID stores a request ID into the given context and returns the new context.
func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// RequestIDMiddleware is a Gin middleware that extracts the request ID from the
// X-Request-Id header. If the header is not present, it generates a new UUID.
// The request ID is stored in the request context and can be retrieved using
// GetRequestID or GetRequestIDFromGin.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx := SetRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)

		// Also set it in the response header for traceability
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}
