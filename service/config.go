package main

import (
	"fmt"
	"os"
)

const (
	envLiteConnectionURL      = "LITE_CONNECTION_URL"
	envTONStakersContractAddr = "TON_STAKERS_ADDRESS"
	envKeyPairSecretName      = "KEY_PAIR_SECRET_NAME"
)

type Config struct {
	Environment string

	LiteConnectionURL         string
	TONStakingContractAddress string
	KeyPairSecretName         string
}

func NewConfig() *Config {
	env := os.Getenv("TARGET_ENV")

	return &Config{
		Environment:               env,
		LiteConnectionURL:         os.Getenv(envLiteConnectionURL),
		TONStakingContractAddress: os.Getenv(envTONStakersContractAddr),
		KeyPairSecretName:         fmt.Sprintf("%s_%s", env, os.Getenv(envKeyPairSecretName)),
	}
}
