package pipes

import (
	"context"
	"net/http"
)

// Middleware function.
type Middleware = func(ctx context.Context, w http.ResponseWriter, r *http.Request, next func(ctx context.Context))

// ContextConstructor context constructor function.
type ContextConstructor = func(w http.ResponseWriter, r *http.Request) context.Context

// Next next middleware function.
type Next = func(context.Context)

// ContextKey pipes context key type.
type ContextKey string

func (p ContextKey) String() string {
	return "Pipes. Context key: " + string(p)
}

// NewPipeline creates new empty pipeline.
func NewPipeline(contextConstructor ContextConstructor) Pipeline {
	p := Pipeline{
		middlewares:        []Middleware{},
		contextConstructor: contextConstructor,
	}
	return p
}

// Pipeline instance.
type Pipeline struct {
	middlewares        []Middleware
	contextConstructor ContextConstructor
}

// Run returns a function which can run the pipeline.
func (p Pipeline) Run() func(w http.ResponseWriter, r *http.Request) {
	var next Next
	return func(w http.ResponseWriter, r *http.Request) {
		middlewareCount := len(p.middlewares)
		if middlewareCount == 0 {
			return
		}

		ctx := p.contextConstructor(w, r)
		i := -1

		next = func(ctx context.Context) {
			i++
			if i <= middlewareCount-1 {
				p.middlewares[i](ctx, w, r, next)
			}
		}
		next(ctx)
	}
}

// Use adds a new middleware to the pipeline and returns a new pipeline.
func (p Pipeline) Use(m Middleware) Pipeline {
	return Pipeline{
		middlewares:        append(p.middlewares, m),
		contextConstructor: p.contextConstructor,
	}
}
