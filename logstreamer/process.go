package logstreamer

import (
	"context"
	"fmt"
)

type Server struct {
	// logstreamer.UnimplementedSimpleServer
}

func (*Server) ProcessRequest(context context.Context, req *LogRequest) (*LogResponse, error) {
	fmt.Println("Got a new Add request")
	mins := req.GetMinutes()
	// profiles := req.GetProfiles()
	result := &LogResponse{Result: mins * 10000}
	return result, nil
}

func (*Server) mustEmbedUnimplementedLogStreamerServer() {

}
