package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/dgtm/log360/logstreamer"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":50551", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := logstreamer.NewLogStreamerClient(conn)

	response, err := c.ProcessRequest(context.Background(), &logstreamer.LogRequest{Minutes: 10})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %d", response.Result)

}
