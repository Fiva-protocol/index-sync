package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

func handler(_ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sigCtx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	config := GetConfig()
	_ = sigCtx

	client := liteclient.NewConnectionPool()
	if err := client.AddConnectionsFromConfigUrl(sigCtx, config.LiteConnectionURL); err != nil {
		fmt.Println("add connections err:", err)
		return InternalServerErr, err
	}

	var (
		ctx            = client.StickyContext(sigCtx)
		tonStakingAddr = address.MustParseAddr(config.TONStakingContractAddress)
	)

	index, err := retryCalculateIndex(ctx, config.LiteConnectionURL, tonStakingAddr, 5, time.Second*2)
	if err != nil {
		fmt.Println("calculate index err:", err)
		return InternalServerErr, err
	}

	var out = ResponseSuccess{
		Index:     index.Uint64(),
		Timestamp: time.Now().UTC().Unix(),
	}

	privateKey, err := GetPrivateKey(ctx, config)
	if err != nil {
		fmt.Println("get privateKey err:", err)
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
		fmt.Println("json marshal err:", err)
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
