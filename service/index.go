package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
)

func calculateIndex(
	ctx context.Context,
	url string,
	addr *address.Address,
) (*big.Int, error) {
	var (
		index  = new(big.Float)
		client = liteclient.NewConnectionPool()
	)

	if err := client.AddConnectionsFromConfigUrl(ctx, url); err != nil {
		log.Fatalln("add connections err:", err)
	}

	api := ton.NewAPIClient(client)

	block, err := api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return nil, err
	}

	res, err := api.WaitForBlock(block.SeqNo).RunGetMethod(ctx, block, addr, "get_pool_full_data")
	if err != nil {
		return nil, err
	}

	tuple := res.AsTuple()
	totalBalance, ok := tuple[2].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("total balance is not int")
	}

	supply, ok := tuple[13].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("supply is not int")
	}

	log.Println("total balance:", totalBalance)
	log.Println("supply:", supply)

	floatTotalBalance := new(big.Float).SetInt(totalBalance)
	floatSupply := new(big.Float).SetInt(supply)

	index = index.Quo(floatTotalBalance, floatSupply)
	index = index.Mul(index, big.NewFloat(1000.0))
	log.Println("index:", index.String())

	i, _ := index.Int64()

	return big.NewInt(i), nil
}

func retryCalculateIndex(
	ctx context.Context,
	liteConnectionsURL string,
	tonStakingAddr *address.Address,
	maxRetries int,
	retryDelay time.Duration,
) (*big.Int, error) {
	var index *big.Int
	var err error

	for i := 0; i < maxRetries; i++ {
		index, err = calculateIndex(ctx, liteConnectionsURL, tonStakingAddr)
		if err == nil {
			return index, nil
		}

		log.Printf("Retry %d/%d: calculateIndex error: %v", i+1, maxRetries, err)
		time.Sleep(retryDelay)
	}

	return nil, err
}
