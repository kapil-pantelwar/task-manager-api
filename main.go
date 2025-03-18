package main

import (
"database/sql"
    "log"
    "os"
    "net/http"
    _ "github.com/lib/pq" // New Postgres driver
    "github.com/go-chi/chi/v5"
    "task-manager/controller"
    "task-manager/repository"
    "task-manager/usecase"
)

func main() {
   // Fetch env vars with defaults
   dbHost := os.Getenv("DB_HOST")
   if dbHost == "" {
       dbHost = "localhost" // Default for local runs
   }
   dbPort := os.Getenv("DB_PORT")
   if dbPort == "" {
       dbPort = "5432"
   }
   dbUser := os.Getenv("DB_USER")
   if dbUser == "" {
       dbUser = "postgres"
   }
   dbPassword := os.Getenv("DB_PASSWORD")
   if dbPassword == "" {
       dbPassword = "password"
   }
   dbName := os.Getenv("DB_NAME")
   if dbName == "" {
       dbName = "taskmanager"
   }

   // Build connection string
   connStr := "host=" + dbHost + " port=" + dbPort + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " sslmode=disable"
   
    db, err := sql.Open("postgres", connStr)
   if err != nil {
       log.Fatal("Failed to open database:", err)
   }
    defer db.Close()
    taskRepo := repository.NewTaskPostgresRepo(db)
    userRepo := repository.NewUserPostgresRepo(db)
    taskUC := usecase.NewTaskUseCase(taskRepo)
    authUC := usecase.NewAuthUseCase(userRepo)
    ctrl := controller.NewTaskController(taskUC, authUC)

    // Create chi router
    r := chi.NewRouter()

    // Middleware stack (applied to all routes)
    r.Use(LoggingMiddleware)
	r.Use(ContentTypeMiddleware)
   
   

    // Routes
    r.Post("/login", ctrl.Login)
    r.Post("/logout", ctrl.Logout)

    // Group for /tasks routes
    r.Route("/tasks", func(r chi.Router) {
		r.Use(AuthMiddleWare(authUC))
		r.Post("/", ctrl.CreateTask)      // POST /tasks
        r.Get("/", ctrl.GetAllTasks)      // GET /tasks
        r.Route("/{id}", func(r chi.Router) {
            r.Get("/", ctrl.GetTaskByID)  // GET /tasks/{id}
            r.Patch("/", ctrl.UpdateTask) // PATCH /tasks/{id}
            r.Delete("/", ctrl.DeleteTask) // DELETE /tasks/{id}
        })
    })

    // Start server
    log.Println("Server starting on :8080...")
    err = http.ListenAndServe(":8080", r)
    if err != nil {
        log.Fatal("Server failed:", err)
    }
}