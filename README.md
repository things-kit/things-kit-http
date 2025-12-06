# Things-Kit HTTP

**HTTP Server Interface for Things-Kit**

This module defines the HTTP server abstraction for Things-Kit applications. It contains **only interfaces**, no implementation.

## Installation

```bash
go get github.com/things-kit/things-kit-http
```

## Purpose

The `things-kit-http` package defines the contract that all HTTP server implementations must follow. This allows applications to program against a stable interface while being free to choose any HTTP framework (Gin, Echo, Chi, stdlib, etc.).

## Interface

```go
type Server interface {
    // GetEngine returns the underlying HTTP router/engine for advanced usage
    GetEngine() any
    
    // RegisterHandler registers a handler function with the given HTTP method and path
    RegisterHandler(method, path string, handler any)
}
```

## Available Implementations

### things-kit-httpgin (Recommended)

The [things-kit-httpgin](https://github.com/things-kit/things-kit-httpgin) module provides a Gin-based implementation.

```go
import (
    "github.com/things-kit/things-kit/app"
    "github.com/things-kit/things-kit-httpgin"
)

func main() {
    app.New(
        viperconfig.Module,
        logging.Module,
        httpgin.Module,  // Provides http.Server
        
        fx.Invoke(RegisterRoutes),
    ).Run()
}
```

## Creating Your Own Implementation

You can create custom HTTP server implementations using any framework:

```go
package myhttp

import "github.com/things-kit/things-kit-http"

type MyHTTPServer struct {
    engine *YourFramework
}

func (s *MyHTTPServer) GetEngine() any {
    return s.engine
}

func (s *MyHTTPServer) RegisterHandler(method, path string, handler any) {
    // Register handler with your framework
}
```

## License

MIT License - see LICENSE file for details
