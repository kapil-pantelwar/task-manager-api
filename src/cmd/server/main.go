package main

import (
    "context"
    "net"
    "log"
    "net/http"
    "task-manager/src/internal/adaptors/persistance" 
    "task-manager/src/internal/interfaces/input/api/rest/handler"
    "task-manager/src/internal/interfaces/input/api/rest/routes"
    "task-manager/src/internal/usecase"
    pb "task-manager/src/proto"
    "google.golang.org/grpc"
)

// gRPC server implementation
type taskServer struct {
    pb.UnimplementedTaskServiceServer
    taskUC *usecase.TaskUseCase
}

func (s *taskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest)(*pb.CreateTaskResponse,error){
    task := req.GetTask()
    createdTask, err := s.taskUC.Create(&task.Title, &task.Description,&task.Status)
    if err != nil {
        return nil, err
    }
    return &pb.CreateTaskResponse{
        Task: &pb.Task{
            Id: int32(createdTask.ID),
            Title: createdTask.Title,
            Description: createdTask.Description,
            Status: createdTask.Status,
        },
        }, nil
    
}

func (s *taskServer) GetTasks(ctx context.Context, req *pb.GetTasksRequest)(*pb.GetTasksResponse, error){
    tasks, err := s.taskUC.GetAll()
    if err != nil {
        return nil, err
    }
    pbTasks := make([]*pb.Task, len(tasks))
    for i,t:= range tasks {
        pbTasks[i] = &pb.Task{
            Id: int32(t.ID),
            Title: t.Title,
            Description: t.Description,
            Status: t.Status,
        }
    }
    return &pb.GetTasksResponse{Tasks: pbTasks},nil
}

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

    // Start REST server in a goroutine
    go func() {
        router := routes.SetupRoutes(ctrl,authUC)
        log.Println("REST server starting on :8080...")
        if err := http.ListenAndServe(":8080", router); err != nil {
            log.Fatalf("REST server failed: %v", err)
        }
    }()

    // Start gRPC server

    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    grpcServer := grpc.NewServer()
    pb.RegisterTaskServiceServer(grpcServer,&taskServer{taskUC: taskUC})
    log.Println("gRPC server starting on :50051...")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("gRPC server failed: %v", err)
    }
}