package main

import (
	"context"
	"log"
	pb "task-manager/src/proto"
	"google.golang.org/grpc"
)

func main(){
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	c:= pb.NewTaskServiceClient(conn)

	resp, err := c.CreateTask(context.Background(),&pb.CreateTaskRequest{Task: &pb.Task{Title: "gRPC Task", Description: "Test", Status: "pending"}})
	if err != nil {
		log.Fatalf("CreateTask failed: %v", err)
	}
	log.Printf("Created: %v", resp.GetTask())
	tasks, err := c.GetTasks(context.Background(),&pb.GetTasksRequest{})
	if err != nil {
		log.Fatalf("GetTasks failed: %v", err)
	}
	log.Printf("Tasks: %v", tasks.GetTasks())
}