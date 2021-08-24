package main

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
	requests := [2]*QueryRequest{inhouserequest, crossRequest}

	wg := sync.WaitGroup{}
	wg.Add(len(requests))

	for _, request := range requests {
		go func(request *QueryRequest) {
			fetchLogByProfile(request, &wg)
		}(request)
	}
	wg.Wait()
}

func fetchLogByProfile(request *QueryRequest, wg *sync.WaitGroup) *cloudwatchlogs.GetQueryResultsOutput {
	defer wg.Done()
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(request.AWSProfile),
	)
	log.Print(request.AWSProfile)

	if err != nil {
		log.Fatal(err)
	}
	client := cloudwatchlogs.NewFromConfig(cfg)
	// from := time.Now().Add(-time.Hour*time.Duration(5)).Unix() * 1000
	//logGroupNames := []string{"/ecs/seizure", "/ecs/seizure-worker"}
	// var logGroupNames []string = []string{"/ecs/seizure", "/ecs/seizure-worker"}

	endsAt := time.Now()
	endTime := endsAt.Unix()
	startTime := endsAt.Add(-time.Hour * 24).Unix()

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

	result_input := &cloudwatchlogs.GetQueryResultsInput{QueryId: q.QueryId}
	resp, err := client.GetQueryResults(context.TODO(), result_input)

	if err != nil {
		log.Fatal("err", err)
	}
	// resultWaitGroup = sync.WaitGroup{}
	// go func(request *QueryRequest) {
	// 	getQueryResults(request, &resultWaitGroup)
	// }(request)

	for resp.Status == "Running" || resp.Status == "Scheduled" {
		// responseChan := make(chan string)
		// resp, err := client.GetQueryResults(context.TODO(), result_input)
		// close(responseChan)
		resp, err := client.GetQueryResults(context.TODO(), result_input)
		log.Print(resp.Status)
		results := resp.Results

		for _, elem := range results {
			// log.Print(elem)
			for _, e := range elem {
				// log.Print(request.AWSProfile)
				log.Print(e)
				// log.Print(*e.Value)
			}
			// log.Print(elem.Value)
		}
		if err != nil {
			log.Fatal(err)
		}
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

	return resp

}
