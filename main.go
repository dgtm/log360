package main

import (
	"log"
	"net"

	"github.com/dgtm/log360/logstreamer"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50551")
	if err != nil {
		log.Fatalf("Failed to listen on port 50551: %v", err)
	}

	logServer := logstreamer.Server{}

	s := grpc.NewServer()
	log.Print(s.GetServiceInfo())
	logstreamer.RegisterLogStreamerServer(s, &logServer)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)
	}

}
