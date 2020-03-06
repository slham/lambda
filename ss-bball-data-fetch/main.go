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
        "gopkg.in/yaml.v2"
        "io/ioutil"
        "net/http"
        "os"
        "time"
)

//Handler to fetch nba player stats for the year 2020 and store as yaml file in S3
func HandleRequest(ctx context.Context) (string, error) {
        bucket, nbaKey, nbaUrl := os.Getenv("BUCKET"), os.Getenv("NBA_API_KEY"), os.Getenv("NBA_API_URL")

        // Fetch player stat data
        res, err := http.Get( nbaUrl + "?key=" + nbaKey)
        if err != nil {
                return "failed to fetch nba player stat data", err
        }

        //Read players json response
        defer res.Body.Close()
        body, err := ioutil.ReadAll(res.Body)

        var players []Player
        err = json.Unmarshal(body, &players)
        if err != nil {
        	return "failed to unmarshal player stat response", err
        }

        //Convert players to yaml
        playersBytes, err := yaml.Marshal(players)
        if err != nil {
                return "failed to marshal player stats to yaml []byte", err
        }

        // The session the S3 Uploader will use
        sess := session.Must(session.NewSession())

        // Create an uploader with the session and default options
        uploader := s3manager.NewUploader(sess)

        key := fmt.Sprintf("player-stats/2020/%d.yaml", time.Now().Unix())

        // Upload the yaml file to S3.
        result, err := uploader.Upload(&s3manager.UploadInput{
                Bucket: aws.String(bucket),
                Key:    aws.String(key),
                Body:   bytes.NewReader(playersBytes),
        })
        if err != nil {
                return "failed to upload file, %v", err
        }

        return fmt.Sprintf("file uploaded to, %s\n", aws.StringValue(&result.Location)), nil
}

func main() {
        lambda.Start(HandleRequest)
}

type Player struct {
        Id              int       `json:"PlayerID" yaml:"id"`
        TeamId          int       `json:"TeamID" yaml:"teamId,omitempty"`
        Name            string    `json:"Name" yaml:"name"`
        Position        string    `json:"Position" yaml:"pos,omitempty"`
        Min             int       `json:"Minutes" yaml:"min,omitempty"`
        Fgm             float32   `json:"FieldGoalsMade" yaml:"fgm,omitempty"`
        Fga             float32   `json:"FieldGoalsAttempted" yaml:"fga,omitempty"`
        Fgp             float32   `json:"FieldGoalsPercentage" yaml:"fgp,omitempty"`
        Ftm             float32   `json:"FreeThrowsMade" yaml:"ftm,omitempty"`
        Fta             float32   `json:"FreeThrowsAttempted" yaml:"fta,omitempty"`
        Ftp             float32   `json:"FreeThrowsPercentage" yaml:"ftp,omitempty"`
        Tpm             float32   `json:"ThreePointersMade" yaml:"tpm,omitempty"`
        Tpa             float32   `json:"ThreePointersAttempted" yaml:"tpa,omitempty"`
        Tpp             float32   `json:"ThreePointersPercentage" yaml:"tpp,omitempty"`
        Reb             float32   `json:"TotalReboundsPercentage" yaml:"reb,omitempty"`
        Ass             float32   `json:"AssistsPercentage" yaml:"ass,omitempty"`
        Stl             float32   `json:"StealsPercentage" yaml:"stl,omitempty"`
        Blk             float32   `json:"BlocksPercentage" yaml:"bks,omitempty"`
        Tvs             float32   `json:"TurnOversPercentage" yaml:"tvs,omitempty"`
        Dds             float32   `json:"DoubleDoubles" yaml:"dds,omitempty"`
        Pts             float32   `json:"Points" yaml:"pts,omitempty"`
        Gms             int       `json:"Games" yaml:"gms,omitempty"`
        CreatedDateTime time.Time `yaml:"createdDateTime"`
        UpdatedDateTime time.Time `yaml:"updatedDateTime"`
        Score           float32   `yaml:"score"`
}
