package main

import (
        "bytes"
        "context"
        "encoding/json"
        "fmt"
        "github.com/aws/aws-lambda-go/lambda"
        "github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/aws/session"
        "github.com/aws/aws-sdk-go/service/s3/s3manager"
        "github.com/google/uuid"
        "os"
)

type MyEvent struct {
        Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
        bucket := os.Getenv("BUCKET")
        sampleJsonString, err := json.Marshal(name)
        if err != nil {
                return "failed to marshal 'name'", err
        }
        // The session the S3 Uploader will use
        sess := session.Must(session.NewSession())

        // Create an uploader with the session and default options
        uploader := s3manager.NewUploader(sess)

        u := uuid.New()
        key := fmt.Sprintf("responses/%s", u.String())

        // Upload the file to S3.
        result, err := uploader.Upload(&s3manager.UploadInput{
                Bucket: aws.String(bucket),
                Key:    aws.String(key),
                Body:   bytes.NewReader(sampleJsonString),
        })
        if err != nil {
                return "failed to upload file, %v", err
        }

        return fmt.Sprintf("file uploaded to, %s\n", aws.StringValue(&result.Location)), nil
}

func main() {
        lambda.Start(HandleRequest)
}
