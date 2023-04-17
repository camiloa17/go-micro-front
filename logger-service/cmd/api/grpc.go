package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"log-service/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()
	log.Println(input)

	// write the log

	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)

	if err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	// return response

	res := &logs.LogResponse{Result: "logged!"}

	return res, nil

}

func (app *Config) gRPCListen() {
	list, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	server := grpc.NewServer()

	logs.RegisterLogServiceServer(server, &LogServer{Models: app.Models})
	log.Printf("gRPC Server started on port %s", grpcPort)

	if err := server.Serve(list); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
