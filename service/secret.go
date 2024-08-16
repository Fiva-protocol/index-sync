package main

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

var secretCache, _ = secretcache.New()

type KeyPair struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

func GetPrivateKey(ctx context.Context, cfg *Config) (ed25519.PrivateKey, error) {
	result, err := secretCache.GetSecretStringWithContext(ctx, cfg.KeyPairSecretName)
	if err != nil {
		log.Fatalln("load aws default config err:", err)
		return nil, err
	}

	var keypair KeyPair
	if err = json.Unmarshal([]byte(result), &keypair); err != nil {
		return nil, err
	}

	return hexStringToEd25519PrivateKey(keypair.PrivateKey)
}

func hexStringToEd25519PrivateKey(hexKey string) (ed25519.PrivateKey, error) {
	privateKeyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %w", err)
	}

	if len(privateKeyBytes) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid private key length: %d", len(privateKeyBytes))
	}

	return privateKeyBytes, nil
}
