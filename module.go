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

// HTTP status codes as registered with IANA.
// See: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
const (
	StatusContinue           = 100 // RFC 9110, 15.2.1
	StatusSwitchingProtocols = 101 // RFC 9110, 15.2.2
	StatusProcessing         = 102 // RFC 2518, 10.1
	StatusEarlyHints         = 103 // RFC 8297

	StatusOK                   = 200 // RFC 9110, 15.3.1
	StatusCreated              = 201 // RFC 9110, 15.3.2
	StatusAccepted             = 202 // RFC 9110, 15.3.3
	StatusNonAuthoritativeInfo = 203 // RFC 9110, 15.3.4
	StatusNoContent            = 204 // RFC 9110, 15.3.5
	StatusResetContent         = 205 // RFC 9110, 15.3.6
	StatusPartialContent       = 206 // RFC 9110, 15.3.7
	StatusMultiStatus          = 207 // RFC 4918, 11.1
	StatusAlreadyReported      = 208 // RFC 5842, 7.1
	StatusIMUsed               = 226 // RFC 3229, 10.4.1

	StatusMultipleChoices   = 300 // RFC 9110, 15.4.1
	StatusMovedPermanently  = 301 // RFC 9110, 15.4.2
	StatusFound             = 302 // RFC 9110, 15.4.3
	StatusSeeOther          = 303 // RFC 9110, 15.4.4
	StatusNotModified       = 304 // RFC 9110, 15.4.5
	StatusUseProxy          = 305 // RFC 9110, 15.4.6
	StatusTemporaryRedirect = 307 // RFC 9110, 15.4.8
	StatusPermanentRedirect = 308 // RFC 9110, 15.4.9

	StatusBadRequest                   = 400 // RFC 9110, 15.5.1
	StatusUnauthorized                 = 401 // RFC 9110, 15.5.2
	StatusPaymentRequired              = 402 // RFC 9110, 15.5.3
	StatusForbidden                    = 403 // RFC 9110, 15.5.4
	StatusNotFound                     = 404 // RFC 9110, 15.5.5
	StatusMethodNotAllowed             = 405 // RFC 9110, 15.5.6
	StatusNotAcceptable                = 406 // RFC 9110, 15.5.7
	StatusProxyAuthRequired            = 407 // RFC 9110, 15.5.8
	StatusRequestTimeout               = 408 // RFC 9110, 15.5.9
	StatusConflict                     = 409 // RFC 9110, 15.5.10
	StatusGone                         = 410 // RFC 9110, 15.5.11
	StatusLengthRequired               = 411 // RFC 9110, 15.5.12
	StatusPreconditionFailed           = 412 // RFC 9110, 15.5.13
	StatusRequestEntityTooLarge        = 413 // RFC 9110, 15.5.14
	StatusRequestURITooLong            = 414 // RFC 9110, 15.5.15
	StatusUnsupportedMediaType         = 415 // RFC 9110, 15.5.16
	StatusRequestedRangeNotSatisfiable = 416 // RFC 9110, 15.5.17
	StatusExpectationFailed            = 417 // RFC 9110, 15.5.18
	StatusTeapot                       = 418 // RFC 9110, 15.5.19 (Unused)
	StatusMisdirectedRequest           = 421 // RFC 9110, 15.5.20
	StatusUnprocessableEntity          = 422 // RFC 9110, 15.5.21
	StatusLocked                       = 423 // RFC 4918, 11.3
	StatusFailedDependency             = 424 // RFC 4918, 11.4
	StatusTooEarly                     = 425 // RFC 8470, 5.2.
	StatusUpgradeRequired              = 426 // RFC 9110, 15.5.22
	StatusPreconditionRequired         = 428 // RFC 6585, 3
	StatusTooManyRequests              = 429 // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  = 431 // RFC 6585, 5
	StatusUnavailableForLegalReasons   = 451 // RFC 7725, 3

	StatusInternalServerError           = 500 // RFC 9110, 15.6.1
	StatusNotImplemented                = 501 // RFC 9110, 15.6.2
	StatusBadGateway                    = 502 // RFC 9110, 15.6.3
	StatusServiceUnavailable            = 503 // RFC 9110, 15.6.4
	StatusGatewayTimeout                = 504 // RFC 9110, 15.6.5
	StatusHTTPVersionNotSupported       = 505 // RFC 9110, 15.6.6
	StatusVariantAlsoNegotiates         = 506 // RFC 2295, 8.1
	StatusInsufficientStorage           = 507 // RFC 4918, 11.5
	StatusLoopDetected                  = 508 // RFC 5842, 7.2
	StatusNotExtended                   = 510 // RFC 2774, 7
	StatusNetworkAuthenticationRequired = 511 // RFC 6585, 6
)
