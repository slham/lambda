   GOOS=linux go build main.go
   zip function.zip main
   aws lambda update-function-code --function-name sheldonsandbox-basketball-data-source --zip-file fileb://function.zip --profile tharivol
