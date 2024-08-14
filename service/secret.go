package main

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

const versionStage = "AWSCURRENT"

func GetPrivateKey(ctx context.Context, cfg *Config) (ed25519.PrivateKey, error) {
	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.AWSConfig.Region))
	if err != nil {
		log.Fatalln("load aws default config err:", err)
		return nil, err
	}

	svc := secretsmanager.NewFromConfig(awsConfig)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(cfg.AWSConfig.SecretManagerSecretName),
		VersionStage: aws.String(versionStage),
	}

	result, err := svc.GetSecretValue(ctx, input)
	if err != nil {
		log.Fatalln("get secret value err:", err)
		return nil, err
	}

	type KeyPair struct {
		PrivateKey string `json:"private_key"`
		PublicKey  string `json:"public_key"`
	}
	var keypair KeyPair
	if err = json.Unmarshal([]byte(*result.SecretString), &keypair); err != nil {
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
