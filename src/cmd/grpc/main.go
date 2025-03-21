package main

import (
	"context"
	"log"
	"net"
	"task-manager/src/internal/adaptors/persistance"
	//"task-manager/src/internal/core/session"
	"task-manager/src/internal/usecase"
	pb "task-manager/src/internal/interfaces/input/grpc/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type taskServer struct {
	pb.UnimplementedTaskServiceServer
	taskUC *usecase.TaskUseCase
	authUC *usecase.AuthUseCase
}

func (s *taskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest)(*pb.CreateTaskResponse, error){
	if err := s.authorize(ctx); err != nil {
		return nil, err
	}
	task := req.GetTask()
	createdTask, err := s.taskUC.Create(&task.Title, &task.Description, &task.Status)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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

func (s *taskServer) GetTasks(req *pb.GetTasksRequest, stream pb.TaskService_GetTasksServer) error {
	if err := s.authorize(stream.Context()); err != nil {
		return err
	}
	tasks, err := s.taskUC.GetAll()
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	for _,t := range tasks {
		if err := stream.Send(&pb.GetTasksResponse{
			Task: &pb.Task{
				Id: int32(t.ID),
				Title: t.Title,
				Description: t.Description,
				Status: t.Status,
			},
		}); err != nil {
			return status.Error(codes.Internal, "failed to send task")
		}
		
	}
	return nil
}

func (s *taskServer) authorize(ctx context.Context) error {
	if s.authUC == nil {
		return status.Error(codes.Internal, "auth service unavailable")
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get("session-id")) == 0 {
		return status.Error(codes.Unauthenticated, "session-id missing")
	}
	sessionID := md.Get("session-id")[0]
	authorized, err := s.authUC.Authorize(sessionID,"user")
	if err != nil || !authorized {
		return status.Error(codes.Unauthenticated,"invalid or unauthorized session")
	}
	return nil
}

func main() {
	//Initialize database
	db, err := persistance.NewDatabase()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	if db.GetDB() == nil {
		log.Fatal("Database connection is nil")
	}

	// Set up use cases
	taskRepo := persistance.NewTaskPostgresRepo(db.GetDB())
    userRepo := persistance.NewUserPostgresRepo(db.GetDB())
    sessionRepo := persistance.NewSessionPostgresRepo(db.GetDB())
    taskUC := usecase.NewTaskUseCase(taskRepo)
    authUC := usecase.NewAuthUseCase(userRepo, sessionRepo)
    if authUC == nil {
        log.Fatal("AuthUseCase is nil")
    }

    // Start gRPC server
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    grpcServer := grpc.NewServer()
    pb.RegisterTaskServiceServer(grpcServer, &taskServer{taskUC: taskUC, authUC: authUC})
    log.Println("gRPC server starting on :50051...")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("gRPC server failed: %v", err)
    }

}