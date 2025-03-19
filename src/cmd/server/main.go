package main

import (
    "log"
    "net/http"
    "task-manager/src/internal/adaptors/persistance" 
    "task-manager/src/internal/interfaces/input/api/rest/handler"
    "task-manager/src/internal/interfaces/input/api/rest/routes"
    "task-manager/src/internal/usecase"
)

func main() {
    // Initialize database
    db, err := persistance.NewDatabase()
    if err != nil {
        log.Fatal("Failed to initialize database:", err)
    }
    defer db.Close()

    // Dependency injection
    taskRepo := persistance.NewTaskPostgresRepo(db.GetDB())
    userRepo := persistance.NewUserPostgresRepo(db.GetDB())
    sessionRepo := persistance.NewSessionPostgresRepo(db.GetDB())
    taskUC := usecase.NewTaskUseCase(taskRepo)
    authUC := usecase.NewAuthUseCase(userRepo, sessionRepo)
    ctrl := controller.NewTaskController(taskUC, authUC)

    // Set up routes
    router := routes.SetupRoutes(ctrl, authUC)

    log.Println("Server starting on :8080...")
    err = http.ListenAndServe(":8080", router)
    if err != nil {
        log.Fatal("Server failed:", err)
    }
}