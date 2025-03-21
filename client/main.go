package main

import (
	"context"
	"io"
	"log"
	pb "task-manager/src/internal/interfaces/input/grpc/task"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main(){
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	c:= pb.NewTaskServiceClient(conn)  //grpc stub

	// sessionID from REST login
	sessionID := "U3KvqpM_ysliquhGyhMG-A==" // tbr with actual session ID from POST /login
	ctx := metadata.AppendToOutgoingContext(context.Background(), "session-id", sessionID)
	
	resp, err:= c.CreateTask(ctx, &pb.CreateTaskRequest{
		Task: &pb.Task{Title: "gRPC Task", Description: "Test 00", Status: "pending"},
	})
	if err != nil {
		log.Fatalf("CreateTask failed: %v", err)
	}

	log.Printf("Created: %v", resp.GetTask())
	stream, err := c.GetTasks(ctx,&pb.GetTasksRequest{})
	if err != nil {
		log.Fatalf("GetTasks failed: %v", err)
	}
	for {
		taskResp, err := stream.Recv()
		if err == io.EOF {
			break // End of stream
		}
		if err != nil {
			log.Fatalf("Failed to receive task: %v", err)
		}
		log.Printf("Tasks: %v", taskResp.GetTask())
		time.Sleep(100*time.Millisecond) //streaming delay
	}
	
}