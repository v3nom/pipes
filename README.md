# Pipes [![Build Status](https://travis-ci.com/v3nom/pipes.svg?branch=master)](https://travis-ci.com/v3nom/pipes)
Minimalist HTTP middleware library for Go web applications running on Google App Engine. Works great with Mux library.

## Why to use middlewares and pipelines?
It helps to reuse request handling code and achieve better separation of concerns.

## Usage

[Documentation](https://godoc.org/github.com/v3nom/pipes)

```go
// Code snippet from a real web application.
apiPipeline := pipes.New().
		Use(middleware.AppEngineContext).   // HTTP request will first get Google App Engine context
		Use(panicMiddleware).               // panic handler will recover, log and present user friendly message if panic is called              
		Use(rateLimitMiddleware).           // rate limitter will make sure that our API endpoint is not overwhelmed
		Use(cookieAuthMiddleware).          // auth cookie will be validated and user object added to the context
        Use(authAPIMiddleware)              // if user is missing, unauthorised API response will be returned without running any business logic

router := mux.NewRouter()
router.HandleFunc("/api/test", apiPipeline.Use(apiHandler).Build())
router.HandleFunc("/api/test2", apiPipeline.Use(api2Handler).Build())
http.Handle("/", router)

func apiHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, next pipes.Next){
    fmt.Fprint(w, "OK")
    next(ctx)
}
func api2Handler(ctx context.Context, w http.ResponseWriter, r *http.Request, next pipes.Next){
    fmt.Fprint(w, "OK2")
    next(ctx)
}
```
