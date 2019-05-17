package pipes

import (
	"context"
	"net/http"
)

// Middleware function.
type Middleware = func(ctx context.Context, w http.ResponseWriter, r *http.Request, next Next)

// Next next middleware function.
type Next = func(context.Context)

// ContextKey pipes context key type.
type ContextKey string

func (p ContextKey) String() string {
	return "Pipes. Context key: " + string(p)
}

// New creates new empty pipeline.
func New() Pipeline {
	p := Pipeline{
		middlewares: []Middleware{},
	}
	return p
}

// Pipeline instance.
type Pipeline struct {
	middlewares []Middleware
}

// Build returns a function which can run the pipeline.
func (p Pipeline) Build() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		middlewareCount := len(p.middlewares)
		if middlewareCount == 0 {
			return
		}

		var next Next
		i := -1

		next = func(ctx context.Context) {
			i++
			if i <= middlewareCount-1 {
				p.middlewares[i](ctx, w, r, next)
			}
		}
		next(nil)
	}
}

// Use adds a new middleware to the pipeline and returns a new pipeline.
func (p Pipeline) Use(m Middleware) Pipeline {
	cop := append(p.middlewares[:0:0], p.middlewares...)
	return Pipeline{
		middlewares: append(cop, m),
	}
}
