package grpc

import (
	"log"
	"net"
	"context"
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

func (s *taskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
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
			Id:          int32(createdTask.ID),
			Title:       createdTask.Title,
			Description: createdTask.Description,
			Status:      createdTask.Status,
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
	for _, t := range tasks {
		if err := stream.Send(&pb.GetTasksResponse{
			Task: &pb.Task{
				Id:          int32(t.ID),
				Title:       t.Title,
				Description: t.Description,
				Status:      t.Status,
			},
		}); err != nil {
			return status.Error(codes.Internal, "failed to send task")
		}
	}
	return nil
}

func (s *taskServer) GetTaskByID(ctx context.Context, req *pb.GetTaskByIDRequest) (*pb.GetTaskByIDResponse, error) {
	if err := s.authorize(ctx); err != nil {
		return nil, err
	}
	task, err := s.taskUC.GetByID(int(req.GetId()))
	if err != nil {
		return nil, status.Error(codes.NotFound, "task not found")
	}
	return &pb.GetTaskByIDResponse{
		Task: &pb.Task{
			Id:          int32(task.ID),
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
		},
	}, nil
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
	authorized, err := s.authUC.Authorize(sessionID, "user")
	if err != nil || !authorized {
		return status.Error(codes.Unauthenticated, "invalid or unauthorized session")
	}
	return nil
}

// StartGRPC launches the gRPC server
func StartGRPC(taskUC *usecase.TaskUseCase, authUC *usecase.AuthUseCase) error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	pb.RegisterTaskServiceServer(grpcServer, &taskServer{taskUC: taskUC, authUC: authUC})
	log.Println("gRPC server running on :50051...")
	return grpcServer.Serve(lis)
}