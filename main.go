package main

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"log"
)

const mnemonic = "prevent sweet hen bird adapt cousin please era carpet assault check sugar lion canal task victory naive fix heavy crucial excuse sense apple sample"

func main() {
	client := liteclient.NewConnectionPool()

	configUrl := "https://ton-blockchain.github.io/testnet-global.config.json"
	err := client.AddConnectionsFromConfigUrl(context.Background(), configUrl)
	if err != nil {
		panic(err)
	}

	api := ton.NewAPIClient(client)
	// bound requests to the same node
	ctx := client.StickyContext(context.Background())

	// we need fresh block info to run get methods
	block, err := api.CurrentMasterchainInfo(ctx)
	if err != nil {
		log.Fatalln("get block err:", err.Error())
		return
	}

	// call method to get seqno of contract
	res, err := api.WaitForBlock(block.SeqNo).RunGetMethod(ctx, block, address.MustParseAddr("EQCnF65lKoXuXxq4XWDVT8OHA6C4XUsrGYjgkCsBpBJKgL51"), "get_index")
	if err != nil {
		log.Fatalln("run get method err:", err.Error())
		return
	}

	index := res.MustInt(0)
	fmt.Println(index)

	op := uint64(0x75706461) // "update_index" в hex
	newIndex := uint64(123)  // Пример нового индекса

	bodyToSign := cell.BeginCell().
		MustStoreUInt(op, 32).
		MustStoreUInt(newIndex, 32)

	privateKey, err := hexStringToEd25519PrivateKey("b7f99270b630cf2cfce343f924a241b48075a0168deb2c62f6602d54f6ab8787758c2f306980f710c63dd8545d132b3172ef5a6229376861d4455e7ca1aed8b5")
	if err != nil {
		panic(err)
	}

	fmt.Println(hex.EncodeToString(privateKey.Public().(ed25519.PublicKey)))

	sign := bodyToSign.EndCell().Sign(privateKey)
	payload := cell.BeginCell().MustStoreBuilder(bodyToSign).MustStoreSlice(sign, 512).EndCell()

	msg := &tlb.ExternalMessage{
		DstAddr: address.MustParseAddr("EQCnF65lKoXuXxq4XWDVT8OHA6C4XUsrGYjgkCsBpBJKgL51"),
		Body:    payload,
	}

	err = api.SendExternalMessage(ctx, msg)
	if err != nil {
		// FYI: it can fail if not enough balance on contract
		log.Printf("send err: %s", err.Error())
		return
	}

	log.Println("External message successfully processed and should be added to blockchain soon!")
	log.Println("Rerun this script in a couple seconds and you should see total and seqno changed.")
}

func hexStringToEd25519PrivateKey(hexKey string) (ed25519.PrivateKey, error) {
	// Decode the hex string to bytes
	privateKeyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %v", err)
	}

	// Ensure the length is correct for Ed25519 private keys
	if len(privateKeyBytes) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid private key length: %v", len(privateKeyBytes))
	}

	// Convert to Ed25519 private key
	privateKey := ed25519.PrivateKey(privateKeyBytes)
	return privateKey, nil
}
