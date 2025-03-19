package routes

import (
	"time"
	"log"
    "net/http"
    "github.com/go-chi/chi/v5"
    "task-manager/src/internal/interfaces/input/api/rest/handler"
    "task-manager/src/internal/usecase"
)

// SetupRoutes configures the router with middleware and endpoints
func SetupRoutes(ctrl *controller.TaskController, authUC *usecase.AuthUseCase) http.Handler {
    r := chi.NewRouter()

    // Global middleware
    r.Use(LoggingMiddleware)
    r.Use(ContentTypeMiddleware)

    // Public routes
    r.Post("/login", ctrl.Login)
    r.Post("/logout", ctrl.Logout)

    // Protected /tasks routes
    r.Route("/tasks", func(r chi.Router) {
        r.Use(AuthMiddleware(authUC))
        r.Post("/", ctrl.CreateTask)
        r.Get("/", ctrl.GetAllTasks)
        r.Route("/{id}", func(r chi.Router) {
            r.Get("/", ctrl.GetTaskByID)
            r.Patch("/", ctrl.UpdateTask)
            r.Delete("/", ctrl.DeleteTask)
        })
    })

    return r
}

// Middleware functions (moved from main.go)
type responseWriterWrapper struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("Started %s %s at %v", r.Method, r.URL.Path, start)
        rw := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
        next.ServeHTTP(rw, r)
        duration := time.Since(start)
        log.Printf("Completed %s %s with status %d in %v", r.Method, r.URL.Path, rw.statusCode, duration)
    })
}

func AuthMiddleware(authUC *usecase.AuthUseCase) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            cookie, err := r.Cookie("session_id")
            if err != nil {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            authorized, err := authUC.Authorize(cookie.Value, "user")
            if err != nil || !authorized {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

func ContentTypeMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}