package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

const errMsgFmt = `{"error":"%s"}`

var headers = map[string]string{
	"Content-Type": "application/json",
}

type ResponseSuccess struct {
	Index     uint64 `json:"index"`
	Timestamp int64  `json:"timestamp"`
	Hash      []byte `json:"hash"`
}

var InternalServerErr = events.APIGatewayProxyResponse{
	StatusCode: http.StatusInternalServerError,
	Headers:    headers,
	Body:       fmt.Sprintf(errMsgFmt, "internal Server Error"),
}
