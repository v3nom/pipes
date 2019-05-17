package pipes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePipeline(t *testing.T) {
	pipeline := NewPipeline()
	pipeline1 := pipeline.Use(middlewareA)
	pipeline2 := pipeline1.Use(middlewareB)

	if len(pipeline.middlewares) != 0 {
		t.Fatal("Expected 0 middleware")
	}
	if len(pipeline1.middlewares) != 1 {
		t.Fatal("Expected 1 middleware")
	}
	if len(pipeline2.middlewares) != 2 {
		t.Fatal("Expected 2 middleware")
	}
}

func TestEmptyPipeline(t *testing.T) {
	pipeline := NewPipeline()

	// Request
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Request handling
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(pipeline.Run())
	handler.ServeHTTP(recorder, req)
}

func TestPipelineNext(t *testing.T) {
	pipeline := NewPipeline().Use(middlewareA).Use(middlewareB)

	// Request
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Request handling
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(pipeline.Run())
	handler.ServeHTTP(recorder, req)

	pipelineHeader := recorder.HeaderMap.Get("pipeline")
	if pipelineHeader != "B" {
		t.Fatalf("Pipeline header. Expected: %v, Actual: %v", "B", pipelineHeader)
	}
}

func TestPipelineWithContext(t *testing.T) {
	pipeline := NewPipeline().
		Use(setContextMiddleware).
		Use(func(ctx context.Context, w http.ResponseWriter, r *http.Request, next Next) {
			if ctx.Value(middlewareID) != "1" {
				t.Fatalf("Pass context. Expected: %v, Acutal: %v", "1", ctx.Value("pipeline"))
			}
			next(ctx)
		})

	// Request
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Request handling
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(pipeline.Run())
	handler.ServeHTTP(recorder, req)

	expectedKeyString := "Pipes. Context key: pipeline"
	if middlewareID.String() != expectedKeyString {
		t.Fatalf("Expected: %v, Actual: %v", expectedKeyString, middlewareID.String())
	}
}

func defaultContextMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request, next func(ctx context.Context)) {
	next(context.TODO())
}

func middlewareA(ctx context.Context, w http.ResponseWriter, r *http.Request, next func(ctx context.Context)) {
	w.Header().Set("pipeline", "A")

	next(ctx)
}

func middlewareB(ctx context.Context, w http.ResponseWriter, r *http.Request, next func(ctx context.Context)) {
	w.Header().Set("pipeline", "B")

	next(ctx)
}

const middlewareID ContextKey = "pipeline"

func setContextMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request, next func(ctx context.Context)) {
	ctx = context.WithValue(ctx, middlewareID, "1")

	next(ctx)
}
