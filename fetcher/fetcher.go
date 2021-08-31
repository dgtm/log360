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
	Frame         int32
}

type QueryResult struct {
	*QueryRequest
	Data []string
}

// Body close???

func Fetch(mins int32, partialResponseChan chan *QueryResult) {
	inhouserequest := &QueryRequest{
		AWSProfile:    "trex-testing",
		LogGroupNames: []string{"/ecs/seizure", "/ecs/seizure-worker"},
		CWQuery:       "fields @message | filter @message LIKE /Worker/ | stats count(*) as workerCount by bin(1m) | sort workerCount desc",
		Limit:         10,
		Frame:         mins,
	}

	crossRequest := &QueryRequest{
		AWSProfile:    "diba-testing",
		LogGroupNames: []string{"/ecs/digital-banking"},
		CWQuery:       "fields @message | filter @message LIKE /Worker/ | stats count(*) as workerCount by bin(1m) | sort workerCount desc",
		Limit:         10,
		Frame:         mins,
	}
	requests := []*QueryRequest{crossRequest, inhouserequest}
	result := []*QueryResult{}

	wg := sync.WaitGroup{}
	wg.Add(len(requests))

	fullResponseChan := make(chan *QueryResult, len(requests))

	for _, request := range requests {
		go func(request *QueryRequest) {
			fullResponseChan <- fetchLogByProfile(request, &wg, partialResponseChan)
		}(request)
	}

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
	// log.Print(time.Duration(request.Frame))
	startTime := endsAt.Add(-time.Hour * 24000).Unix()

	input := &cloudwatchlogs.StartQueryInput{
		EndTime:       aws.Int64(endTime),
		StartTime:     aws.Int64(startTime),
		LogGroupNames: request.LogGroupNames,
		QueryString:   aws.String(request.CWQuery),
		Limit:         aws.Int32(10000),
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
		log.Printf("%s %s", request.AWSProfile, resp.Status)
		// log.Print(*result_input.QueryId)
		resp = getResponseStatus(client, result_input)

		resultData := []string{}
		results := resp.Results
		// log.Printf("%+v", resp)

		for _, elem := range results {
			for _, e := range elem {
				log.Printf("%+v", e)
				resultData = append(resultData, *e.Value)
				// log.Printf("%+v", resultData)
				// log.Print(*e.Value)
			}
		}

		result := &QueryResult{
			QueryRequest: request,
			Data:         resultData,
		}
		// log.Print(result)
		partialResponseChan <- result
		time.Sleep(8 * time.Second)

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
		Data:         []string{"asdasd"},
	}
	// log.Print(result)

	return result
}
