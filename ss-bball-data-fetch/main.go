package main

import (
        "bytes"
        "context"
        "fmt"
        "github.com/aws/aws-lambda-go/lambda"
        "github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/aws/session"
        "github.com/aws/aws-sdk-go/service/s3/s3manager"
        "io/ioutil"
        "net/http"
        "os"
        "time"
)

type MyEvent struct {
        Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
        bucket, nbaKey, nbaUrl := os.Getenv("BUCKET"), os.Getenv("NBA_API_KEY"), os.Getenv("NBA_API_URL")

        // Fetch player stat data
        res, err := http.Get( nbaUrl + "?key=" + nbaKey)
        if err != nil {
                return "failed to fetch nba player stat data", err
        }

        defer res.Body.Close()
        body, err := ioutil.ReadAll(res.Body)

        // The session the S3 Uploader will use
        sess := session.Must(session.NewSession())

        // Create an uploader with the session and default options
        uploader := s3manager.NewUploader(sess)

        key := fmt.Sprintf("player-stats/2020/%d.json", time.Now().Unix())

        // Upload the file to S3.
        result, err := uploader.Upload(&s3manager.UploadInput{
                Bucket: aws.String(bucket),
                Key:    aws.String(key),
                Body:   bytes.NewReader(body),
        })
        if err != nil {
                return "failed to upload file, %v", err
        }

        return fmt.Sprintf("file uploaded to, %s\n", aws.StringValue(&result.Location)), nil
}

func main() {
        lambda.Start(HandleRequest)
}
