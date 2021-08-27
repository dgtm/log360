package logstreamer

import (
	"context"
	"fmt"
)

type server struct{}

func (*server) ProcessRequest(context context.Context, req *logstreamer.LogRequest) (*logstreamer.LogResponse, error) {
	fmt.Println("Got a new Add request")
	mins := req.GetMinutes()
	// profiles := req.GetProfiles()
	result := &logstreamer.LogResponse{Result: mins * 100}
	return result, nil
}