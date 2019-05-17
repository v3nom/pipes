package pipes

import (
	"context"
	"net/http"
)

// Middleware function.
type Middleware = func(ctx context.Context, w http.ResponseWriter, r *http.Request, next func(ctx context.Context))

// Next next middleware function.
type Next = func(context.Context)

// ContextKey pipes context key type.
type ContextKey string

func (p ContextKey) String() string {
	return "Pipes. Context key: " + string(p)
}

// NewPipeline creates new empty pipeline.
func NewPipeline() Pipeline {
	p := Pipeline{
		middlewares: []Middleware{},
	}
	return p
}

// Pipeline instance.
type Pipeline struct {
	middlewares []Middleware
}

// Run returns a function which can run the pipeline.
func (p Pipeline) Run() func(w http.ResponseWriter, r *http.Request) {
	var next Next
	return func(w http.ResponseWriter, r *http.Request) {
		middlewareCount := len(p.middlewares)
		if middlewareCount == 0 {
			return
		}

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
	return Pipeline{
		middlewares: append(p.middlewares, m),
	}
}
