package pipes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePipeline(t *testing.T) {
	pipeline := createPipeline()
	pipeline1 := pipeline.Use(pipelineA)
	pipeline2 := pipeline1.Use(pipelineB)

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
	pipeline := createPipeline()

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
	pipeline := createPipeline().Use(pipelineA).Use(pipelineB)

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

func createPipeline() Pipeline {
	return NewPipeline(func(w http.ResponseWriter, r *http.Request) context.Context {
		return context.TODO()
	})
}

func pipelineA(ctx context.Context, w http.ResponseWriter, r *http.Request, next func(ctx context.Context)) {
	w.Header().Set("pipeline", "A")

	next(ctx)
}

func pipelineB(ctx context.Context, w http.ResponseWriter, r *http.Request, next func(ctx context.Context)) {
	w.Header().Set("pipeline", "B")

	next(ctx)
}
