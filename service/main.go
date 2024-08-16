package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

var (
	client *liteclient.ConnectionPool
	config *Config

	once sync.Once
)

func handler(_ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sigCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := client.AddConnectionsFromConfigUrl(sigCtx, config.LiteConnectionURL); err != nil {
		log.Println("add connections err:", err)
		return InternalServerErr, err
	}

	var (
		ctx            = client.StickyContext(sigCtx)
		tonStakingAddr = address.MustParseAddr(config.TONStakingContractAddress)
	)

	index, err := retryCalculateIndex(ctx, config.LiteConnectionURL, tonStakingAddr, 5, time.Second*2)
	if err != nil {
		log.Println("calculate index err:", err)
		return InternalServerErr, err
	}

	var out = ResponseSuccess{
		Index:     index.Uint64(),
		Timestamp: time.Now().UTC().Unix(),
	}

	privateKey, err := GetPrivateKey(ctx, config)
	if err != nil {
		log.Println("get privateKey err:", err)
		return InternalServerErr, err
	}

	sign := cell.BeginCell().
		MustStoreUInt(uint64(out.Timestamp), 64).
		MustStoreUInt(out.Index, 64).
		EndCell().
		Sign(privateKey)

	out.Hash = sign

	body, err := json.Marshal(&out)
	if err != nil {
		log.Println("json marshal err:", err)
		return InternalServerErr, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(handler)
}

func init() {
	once.Do(func() {
		config = NewConfig()
		client = liteclient.NewConnectionPool()
	})
}
