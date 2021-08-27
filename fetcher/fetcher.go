package fetcher

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

type QueryRequest struct {
	Id            string
	LogGroupNames []string
	CWQuery       string
	AWSProfile    string
	Limit         int32
}

type QueryResult struct {
	*QueryRequest
	data *cloudwatchlogs.GetQueryResultsOutput
}

// Body close???

func main() {
	inhouserequest := &QueryRequest{
		AWSProfile:    "trex-testing",
		LogGroupNames: []string{"/ecs/seizure", "/ecs/seizure-worker"},
		CWQuery:       "fields @message",
		Limit:         10,
	}

	crossRequest := &QueryRequest{
		AWSProfile:    "diba-testing",
		LogGroupNames: []string{"/ecs/digital-banking"},
		CWQuery:       "fields @message",
		Limit:         10,
	}
	requests := []*QueryRequest{crossRequest, inhouserequest}
	result := []*QueryResult{}

	wg := sync.WaitGroup{}
	wg.Add(len(requests))

	fullResponseChan := make(chan *QueryResult, len(requests))
	partialResponseChan := make(chan *QueryResult)

	for _, request := range requests {
		go func(request *QueryRequest) {
			fullResponseChan <- fetchLogByProfile(request, &wg, partialResponseChan)
		}(request)
	}

	go func() {
		for results := range partialResponseChan {
			if results != nil {
				log.Print("from partial")
				log.Printf("%+v", results)
			}

		}
	}()

	wg.Wait()
	close(fullResponseChan)
	close(partialResponseChan)

	for results := range fullResponseChan {
		result = append(result, results)
	}

	// log.Print(result)
}
func makeRequest(client *cloudwatchlogs.Client, request *QueryRequest) {
	endsAt := time.Now()
	endTime := endsAt.Unix()
	startTime := endsAt.Add(-time.Minute * 24).Unix()

	input := &cloudwatchlogs.StartQueryInput{
		EndTime:       aws.Int64(endTime),
		StartTime:     aws.Int64(startTime),
		LogGroupNames: request.LogGroupNames,
		QueryString:   aws.String(request.CWQuery),
		Limit:         aws.Int32(request.Limit),
	}
	q, err := client.StartQuery(context.TODO(), input)

	if err != nil {
		log.Fatal(err)
	}
	request.Id = aws.ToString(q.QueryId)
	log.Printf("QueryId=%s", request.Id)
}

func fetchLogByProfile(request *QueryRequest, wg *sync.WaitGroup, partialResponseChan chan *QueryResult) *QueryResult {
	defer wg.Done()
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(request.AWSProfile),
	)
	log.Print(request.AWSProfile)

	if err != nil {
		log.Fatal(err)
	}
	client := cloudwatchlogs.NewFromConfig(cfg)
	makeRequest(client, request)
	result := fetchResultsByRequest(client, request, partialResponseChan)
	return result
}

func getResponseStatus(client *cloudwatchlogs.Client, input *cloudwatchlogs.GetQueryResultsInput) *cloudwatchlogs.GetQueryResultsOutput {
	resp, _ := client.GetQueryResults(context.TODO(), input)
	return resp
}

//protoc logstreamer.proto --go-grpc_out=../server/logstreamerpb/

func fetchResultsByRequest(client *cloudwatchlogs.Client, request *QueryRequest, partialResponseChan chan *QueryResult) *QueryResult {
	result_input := &cloudwatchlogs.GetQueryResultsInput{QueryId: &request.Id}
	resp, err := client.GetQueryResults(context.TODO(), result_input)
	if err != nil {
		log.Fatal("err", err)
	}

	status := resp.Status
	for status == "Running" || status == "Scheduled" {
		log.Print(resp.Status)
		time.Sleep(2 * time.Second)
		log.Print(*result_input.QueryId)

		resp = getResponseStatus(client, result_input)
		result := &QueryResult{
			QueryRequest: request,
			data:         resp,
		}
		partialResponseChan <- result
		status = resp.Status
		if status == "Complete" {
			results := resp.Results
			for _, elem := range results {
				for _, e := range elem {
					log.Print(*e.Value)
					// log.Print(*e.Value)
				}
			}
		}

		// if err != nil {
		// 	log.Fatal(err)
		// }
		// } else {
		// log.Print(resp.Results)
		// break
		// }
	}
	result := &QueryResult{
		QueryRequest: request,
		data:         resp,
	}
	log.Print(result)

	return result
}
