package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"net/http"
	"os"
	"time"
)

func handler(ctx context.Context, s3Event events.S3Event) (string, error) {
	var bucket, key string
	for _, record := range s3Event.Records {
		bucket = record.S3.Bucket.Name
		key = record.S3.Object.Key
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	}))
	// Create an downloader with the session and default options
	downloader := s3manager.NewDownloader(sess)
	log.Println("downloader created")

	// Store object in buffer
	buff := &aws.WriteAtBuffer{}
	_, err := downloader.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "error downloading object", err
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10,
		},
	}

	reader := bytes.NewBuffer(buff.Bytes())
	req, err := http.NewRequest("PUT", os.Getenv("BASKETBALL_URL"), reader)
	if err != nil {
		return "error creating player stats update request", err
	}

	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Accept-Encoding", "gzip,deflate")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil || res.StatusCode != 200{
		return "error sending player stats request", err
	}

	return fmt.Sprintf("successfully sourced %v", key), nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
