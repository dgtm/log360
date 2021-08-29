package logstreamer

import (
	"context"
	"fmt"
)

// go get -u google.golang.org/protobuf/cmd/protoc-gen-go
// go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
// protoc logstreamer.proto --go-grpc_out=./
// protoc logstreamer.proto --go_out=./

// protoc logstreamer.proto --js_out=import_style=commonjs,binary:../webclient --grpc-web_out=import_style=commonjs,mode=grpcwebtext:../webclient
// protoc logstreamer.proto --js_out=library=logstreamer_client,binary:
// protoc logstreamer.proto --js_out=import_style=commonjs,binary:.
//protoc logstreamer.proto --grpc-web_out=import_style=commonjs,mode=grpcwebtext:.

type Server struct {
	// logstreamer.UnimplementedSimpleServer
}

func (*Server) ProcessRequest(context context.Context, req *LogRequest) (*LogResponse, error) {
	fmt.Println("Got a new Add request")
	mins := req.GetMinutes()
	// profiles := req.GetProfiles()
	result := &LogResponse{Result: mins * 100}
	return result, nil
}

func (*Server) mustEmbedUnimplementedLogStreamerServer() {

}
