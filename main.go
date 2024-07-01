package main

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"log"
	"math/big"
	"time"

	cfg "github.com/igefined/go-kit/config"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

func main() {
	sigCtx := cfg.SigTermIntCtx()

	config, err := NewConfig()
	if err != nil {
		log.Fatalln("init config err:", err)
	}

	client := liteclient.NewConnectionPool()

	if err = client.AddConnectionsFromConfigUrl(sigCtx, config.LiteConnectionsURLTestnet); err != nil {
		log.Fatalln("add connections err:", err)
	}

	var (
		api            = ton.NewAPIClient(client)
		ctx            = client.StickyContext(sigCtx)
		masterAddr     = address.MustParseAddr(config.MasterContractAddress)
		tonStakingAddr = address.MustParseAddr(config.TONStakingContractAddress)
	)

	index, err := getIndex(ctx, api, masterAddr)
	if err != nil {
		log.Fatalln("get index err:", err)
	}

	log.Println("previous index:", index)

	index, err = retryCalculateIndex(ctx, config.LiteConnectionsURLMainnet, tonStakingAddr, 5, time.Second*2)
	if err != nil {
		log.Fatalln("calculate index err:", err)
	}

	privateKey, err := GetPrivateKey(ctx, config)
	if err != nil {
		log.Fatalln("get privateKey err:", err)
	}

	if err = updateIndex(ctx, api, privateKey, masterAddr, index.Uint64()); err != nil {
		log.Fatalln("update index err:", err)
	}

	index, err = getIndex(ctx, api, masterAddr)
	if err != nil {
		log.Fatalln("get index err:", err)
	}

	log.Println("current index:", index)
	log.Println("External message successfully processed and should be added to blockchain soon")
}

func getIndex(ctx context.Context, api *ton.APIClient, addr *address.Address) (*big.Int, error) {
	block, err := api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return nil, err
	}

	res, err := api.WaitForBlock(block.SeqNo).RunGetMethod(ctx, block, addr, "get_index")
	if err != nil {
		return nil, err
	}

	return res.MustInt(0), nil
}

func updateIndex(
	ctx context.Context,
	api *ton.APIClient,
	privateKey ed25519.PrivateKey,
	addr *address.Address,
	index uint64,
) error {
	op := uint64(0xe0178583) // "update_index" into hex

	bodyToSign := cell.BeginCell().
		MustStoreUInt(op, 32).
		MustStoreUInt(index, 32)

	sign := bodyToSign.EndCell().Sign(privateKey)
	payload := cell.BeginCell().MustStoreSlice(sign, 512).MustStoreBuilder(bodyToSign).EndCell()

	msg := &tlb.ExternalMessage{
		DstAddr: addr,
		Body:    payload,
	}

	err := api.SendExternalMessage(ctx, msg)
	if err != nil {
		return err
	}

	return nil
}

func calculateIndex(
	ctx context.Context,
	url string,
	addr *address.Address,
) (*big.Int, error) {
	var (
		index = new(big.Float)

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
	liteConnectionsURLMainnet string,
	tonStakingAddr *address.Address,
	maxRetries int,
	retryDelay time.Duration,
) (*big.Int, error) {
	var index *big.Int
	var err error

	for i := 0; i < maxRetries; i++ {
		index, err = calculateIndex(ctx, liteConnectionsURLMainnet, tonStakingAddr)
		if err == nil {
			return index, nil
		}

		log.Printf("Retry %d/%d: calculateIndex error: %v", i+1, maxRetries, err)
		time.Sleep(retryDelay)
	}

	return nil, err
}
