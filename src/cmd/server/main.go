package main

import (
	"log"
	"net/http"
	"task-manager/src/internal/adaptors/persistance"
	controller "task-manager/src/internal/interfaces/input/api/rest/handler"
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

	// Start REST server 

	router := routes.SetupRoutes(ctrl, authUC)
	log.Println("REST server starting on :8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("REST server failed: %v", err)
	}

}
