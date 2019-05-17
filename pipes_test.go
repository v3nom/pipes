package pipes

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePipeline(t *testing.T) {
	pipeline := New()
	pipeline1 := pipeline.Use(middlewareA)
	pipeline2 := pipeline1.Use(middlewareB)
	pipeline3 := pipeline.Use(defaultContextMiddleware)

	if len(pipeline.middlewares) != 0 {
		t.Fatal("Expected 0 middlewares")
	}
	if len(pipeline1.middlewares) != 1 {
		t.Fatal("Expected 1 middleware")
	}
	if len(pipeline2.middlewares) != 2 {
		t.Fatal("Expected 2 middlewares")
	}
	if len(pipeline3.middlewares) != 1 {
		t.Fatal("Expected 1 middleware")
	}
}

func TestEmptyPipeline(t *testing.T) {
	pipeline := New()

	// Request
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Request handling
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(pipeline.Build())
	handler.ServeHTTP(recorder, req)
}

func TestPipelineNext(t *testing.T) {
	pipeline := New().Use(middlewareA).Use(middlewareB)

	// Request
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Request handling
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(pipeline.Build())
	handler.ServeHTTP(recorder, req)

	pipelineHeader := recorder.HeaderMap.Get("pipeline")
	if pipelineHeader != "B" {
		t.Fatalf("Pipeline header. Expected: %v, Actual: %v", "B", pipelineHeader)
	}
}

func TestPipelineWithContext(t *testing.T) {
	pipeline := New().
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
	handler := http.HandlerFunc(pipeline.Build())
	handler.ServeHTTP(recorder, req)

	expectedKeyString := "Pipes. Context key: pipeline"
	if middlewareID.String() != expectedKeyString {
		t.Fatalf("Expected: %v, Actual: %v", expectedKeyString, middlewareID.String())
	}
}

func TestRunPipelineTwice(t *testing.T) {
	pipeline := New().
		Use(defaultContextMiddleware).
		Use(middlewareA).
		Use(middlewareB).
		Use(middlewareA).
		Use(middlewareB)

	// Request
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Requests
	recorder1 := httptest.NewRecorder()
	handler1 := http.HandlerFunc(pipeline.Use(contentAMiddleware).Build())
	recorder2 := httptest.NewRecorder()
	handler2 := http.HandlerFunc(pipeline.Use(contentBMiddleware).Build())

	handler1.ServeHTTP(recorder1, req)
	if recorder1.Code != http.StatusOK || recorder1.Body.String() != "a" {
		t.Fatalf("Expected: a, Actual: %v", recorder1.Body.String())
	}

	handler2.ServeHTTP(recorder2, req)
	if recorder2.Code != http.StatusOK || recorder2.Body.String() != "b" {
		t.Fatalf("Expected: b, Actual: %v", recorder2.Body.String())
	}
}

func defaultContextMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request, next Next) {
	next(context.TODO())
}

func middlewareA(ctx context.Context, w http.ResponseWriter, r *http.Request, next Next) {
	w.Header().Set("pipeline", "A")

	next(ctx)
}

func middlewareB(ctx context.Context, w http.ResponseWriter, r *http.Request, next Next) {
	w.Header().Set("pipeline", "B")

	next(ctx)
}

const middlewareID ContextKey = "pipeline"

func setContextMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request, next Next) {
	ctx = context.WithValue(ctx, middlewareID, "1")

	next(ctx)
}

func statusOKMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request, next Next) {
	w.WriteHeader(http.StatusOK)
	next(ctx)
}

func contentAMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request, next Next) {
	fmt.Fprint(w, "a")

	next(ctx)
}

func contentBMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request, next Next) {
	fmt.Fprint(w, "b")
	next(ctx)
}
