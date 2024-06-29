package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"

	cfg "github.com/igefined/go-kit/config"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

func main() {
	ctx := cfg.SigTermIntCtx()

	config, err := NewConfig()
	if err != nil {
		log.Fatalln("init config err:", err)
	}

	client := liteclient.NewConnectionPool()

	if err = client.AddConnectionsFromConfigUrl(ctx, config.LiteConnectionsURL); err != nil {
		log.Fatalln("add connections err:", err)
	}

	api := ton.NewAPIClient(client)
	ctx = client.StickyContext(ctx)

	block, err := api.CurrentMasterchainInfo(ctx)
	if err != nil {
		log.Fatalln("get block err:", err.Error())
	}

	res, err := api.
		WaitForBlock(block.SeqNo).
		RunGetMethod(ctx, block, address.MustParseAddr(config.MasterContractAddress), "get_index")
	if err != nil {
		log.Fatalln("run get method err:", err.Error())
	}

	previousIndex := res.MustInt(0)
	log.Println("previous index:", previousIndex)

	op := uint64(0x75706461) // "update_index" в hex
	newIndex := uint64(123)  // Пример нового индекса

	bodyToSign := cell.BeginCell().
		MustStoreUInt(op, 32).
		MustStoreUInt(newIndex, 32)

	privateKey, err := hexStringToEd25519PrivateKey("")
	if err != nil {
		panic(err)
	}

	sign := bodyToSign.EndCell().Sign(privateKey)
	payload := cell.BeginCell().MustStoreBuilder(bodyToSign).MustStoreSlice(sign, 512).EndCell()

	msg := &tlb.ExternalMessage{
		DstAddr: address.MustParseAddr(config.MasterContractAddress),
		Body:    payload,
	}

	err = api.SendExternalMessage(ctx, msg)
	if err != nil {
		log.Printf("send external message err: %s", err.Error())
		return
	}

	log.Println("External message successfully processed and should be added to blockchain soon!")
}

func hexStringToEd25519PrivateKey(hexKey string) (ed25519.PrivateKey, error) {
	// Decode the hex string to bytes
	privateKeyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %w", err)
	}

	// Ensure the length is correct for Ed25519 private keys
	if len(privateKeyBytes) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid private key length: %d", len(privateKeyBytes))
	}

	// Convert to Ed25519 private key
	privateKey := ed25519.PrivateKey(privateKeyBytes)
	return privateKey, nil
}
