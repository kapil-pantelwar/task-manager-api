package routes

import (
	
    "net/http"
    "github.com/go-chi/chi/v5"
    "task-manager/src/internal/interfaces/input/api/rest/handler"
    "task-manager/src/internal/interfaces/input/api/rest/middleware"
    "task-manager/src/internal/usecase"
)

// SetupRoutes configures the router with middleware and endpoints
func SetupRoutes(ctrl *controller.TaskController, authUC *usecase.AuthUseCase) http.Handler {
    r := chi.NewRouter()

    // Global middleware
    r.Use(middleware.LoggingMiddleware)
    r.Use(middleware.ContentTypeMiddleware)

    // Public routes
    r.Post("/login", ctrl.Login)
    r.Post("/logout", ctrl.Logout)

    // Protected /tasks routes
    r.Route("/tasks", func(r chi.Router) {
        r.Use(middleware.AuthMiddleWare(authUC))
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