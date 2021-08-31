package logstreamer

import (
	"fmt"
	"log"

	"github.com/dgtm/log360/fetcher"
)

// go get -u google.golang.org/protobuf/cmd/protoc-gen-go
// go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
// protoc logstreamer.proto --go-grpc_out=./
// protoc logstreamer.proto --go_out=./

// protoc logstreamer.proto --js_out=import_style=commonjs:. --grpc-web_out=import_style=commonjs,mode=grpcwebtext:.

type Server struct {
	// logstreamer.UnimplementedSimpleServer
}

func (*Server) ProcessRequest(req *LogRequest, stream LogStreamer_ProcessRequestServer) error {
	fmt.Println("Got a new Add request")
	mins := req.GetMinutes()
	// profiles := req.GetProfiles()
	partialResponseChan := make(chan *fetcher.QueryResult)

	go func() {
		log.Print("start waiting for response")

		for result := range partialResponseChan {
			if result != nil {
				// log.Print("enter waiting for response")
				// log.Print(result.AWSProfile)
				// log.Printf("data is %s", result.Data)

				stream.Send(&LogResponse{Profile: result.AWSProfile, Result: result.Data})
				// log.Print("from partial")
				// log.Printf("%+v", result)
			}

		}
	}()
	fetcher.Fetch(mins, partialResponseChan)

	return nil
	// result := &LogResponse{Result: mins * 100}
	// return result, nil
}

func (*Server) mustEmbedUnimplementedLogStreamerServer() {

}
