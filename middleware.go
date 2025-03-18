package main

import (
	"log"
	"net/http"
	"task-manager/usecase"
	"time"
)

// responseWriterWrapper captures the status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int){
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware logs details of each request
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("Started %s %s at %v", r.Method, r.URL.Path, start)
        
		// Wrap the ResponseWrite to captus status
		rw := &responseWriterWrapper{ResponseWriter: w,statusCode: http.StatusOK}

        // Call the next handler in the chain
        next.ServeHTTP(rw, r)
        
        // Log the duration after the request is handled
        duration := time.Since(start)
        log.Printf("Completed %s %s with status %d in %v", r.Method, r.URL.Path,rw.statusCode, duration)
    })
}
// ContentTypeMiddleware sets the Content-Type header to application/json

func ContentTypeMiddleware(next http.Handler) http.Handler {
   
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){ 
        log.Printf("Setting Content-Type: application/json")
		w.Header().Set("Content-Type","application/json")
		next.ServeHTTP(w,r)
	})
}

//AuthMiddleWare checks for a valid session cookie 
func AuthMiddleWare(authUC *usecase.AuthUseCase) func(http.Handler)http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
            cookie,err := r.Cookie("session_id")
            if err != nil {
                http.Error(w,"Unauthorized", http.StatusUnauthorized) //debug 
                return 
            }
            authorized, err := authUC.Authorize(cookie.Value,"user")// Minimum "user" role
            if err != nil || !authorized {
                http.Error(w,"Unauthorized", http.StatusUnauthorized)
                return 
            }
            next.ServeHTTP(w,r)
        })
    }
}