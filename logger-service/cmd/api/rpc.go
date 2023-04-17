package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"time"
)

// RPCServer is the type for our RPC server. Methods that take this as a receiver are available
// over RPC, as long as they are exported.
type RPCServer struct{}

// RPC data we receive from RPC
type RPCPayload struct {
	Name string
	Data string
}

// Writes our payload to mongo
func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})

	if err != nil {
		log.Println("Error writing to mongo", err)
		return err
	}

	*resp = fmt.Sprintf("Processed payload via RPC: %s", payload.Name)

	return nil

}
