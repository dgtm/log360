package logstreamer

import (
	"context"
	"fmt"
)

type Server struct{}

func (*Server) ProcessRequest(context context.Context, req *LogRequest) (*LogResponse, error) {
	fmt.Println("Got a new Add request")
	mins := req.GetMinutes()
	// profiles := req.GetProfiles()
	result := &LogResponse{Result: mins * 100}
	return result, nil
}
