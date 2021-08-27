package logstreamer

import (
	"context"
	"fmt"

	"github.com/dgtm/log360/logstreamer"
)

type Server struct{}

func (*Server) ProcessRequest(context context.Context, req *logstreamer.LogRequest) (*LogResponse, error) {
	fmt.Println("Got a new Add request")
	mins := req.GetMinutes()
	// profiles := req.GetProfiles()
	result := &logstreamer.LogResponse{Result: mins * 100}
	return result, nil
}
