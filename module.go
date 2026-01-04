// Package http defines framework-level HTTP server abstractions.
// This package provides interfaces that HTTP implementations must satisfy,
// allowing users to swap HTTP frameworks (Gin, Chi, Echo, stdlib, etc.) while
// maintaining compatibility with the framework.
//
// For a production-ready implementation, see the httpgin package.
package http

import (
	"context"
	"io"
	"net/http"
)

// Server represents an HTTP server that can be started and stopped.
// Implementations should handle the lifecycle of the HTTP server including
// graceful shutdown and proper error handling.
type Server interface {
	// Start begins listening for HTTP requests.
	// This should be non-blocking and return immediately after the server starts.
	// The implementation should start a goroutine for the actual serving.
	Start(ctx context.Context) error

	// Stop gracefully shuts down the HTTP server.
	// It should wait for in-flight requests to complete within the context deadline.
	// Implementations should respect the context's cancellation/timeout.
	Stop(ctx context.Context) error

	// Addr returns the address the server is listening on (e.g., ":8080").
	Addr() string
}

// Context represents an HTTP request/response context abstraction.
// This interface provides a framework-agnostic way to interact with HTTP requests and responses.
// Implementations should wrap framework-specific context types (e.g., gin.Context, echo.Context).
type Context interface {
	// Request returns the underlying *http.Request
	Request() *http.Request

	// Context returns the request's context for cancellation and deadlines
	Context() context.Context

	// Param retrieves a URL path parameter by name
	Param(name string) string

	// Query retrieves a URL query parameter by name
	Query(name string) string

	// QueryDefault retrieves a URL query parameter with a default value
	QueryDefault(name, defaultValue string) string

	// GetHeader retrieves a request header by name
	GetHeader(name string) string

	// SetHeader sets a response header
	SetHeader(name, value string)

	// BindJSON binds the request body as JSON to the provided struct
	BindJSON(obj interface{}) error

	// Bind binds the request body to the provided struct (supports multiple formats)
	Bind(obj interface{}) error

	// JSON sends a JSON response with the given status code
	JSON(code int, obj interface{}) error

	// String sends a string response with the given status code
	String(code int, s string) error

	// Status sets the HTTP response status code
	Status(code int)

	// Writer returns the response writer
	Writer() io.Writer
}

// HandlerFunc defines the function signature for HTTP handlers using the abstract Context.
// This is similar to echo.HandlerFunc but framework-agnostic.
type HandlerFunc func(Context) error

// Router represents an abstract HTTP router for registering routes.
// Implementations should wrap framework-specific routers (e.g., gin.Engine, echo.Echo).
type Router interface {
	// GET registers a GET route
	GET(path string, handler HandlerFunc)

	// POST registers a POST route
	POST(path string, handler HandlerFunc)

	// PUT registers a PUT route
	PUT(path string, handler HandlerFunc)

	// DELETE registers a DELETE route
	DELETE(path string, handler HandlerFunc)

	// PATCH registers a PATCH route
	PATCH(path string, handler HandlerFunc)

	// Group creates a route group with the given prefix
	Group(prefix string) Router
}

// Handler represents a component that can register HTTP routes.
// Handlers should use the abstract Router interface to register routes,
// making them framework-agnostic.
type Handler interface {
	// RegisterRoutes registers this handler's routes with the HTTP router.
	RegisterRoutes(router Router)
}

// Config holds common HTTP server configuration.
// Specific implementations may embed this struct and add framework-specific fields.
type Config struct {
	Port int    `mapstructure:"port"` // Port to listen on
	Host string `mapstructure:"host"` // Host to bind to (empty = all interfaces)
}
