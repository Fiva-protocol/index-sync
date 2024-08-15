package main

import (
	"fmt"
	"os"
	"sync"
)

const (
	envLiteConnectionURL      = "LITE_CONNECTION_URL"
	envTONStakersContractAddr = "TON_STAKERS_ADDRESS"
	envKeyPairSecretName      = "KEY_PAIR_SECRET_NAME"
)

var (
	cfg  *Config
	once sync.Once
)

type Config struct {
	Environment string

	LiteConnectionURL         string
	TONStakingContractAddress string
	KeyPairSecretName         string
}

func GetConfig() *Config {
	once.Do(func() {
		env := os.Getenv("TARGET_ENV")

		cfg = &Config{
			Environment:               env,
			LiteConnectionURL:         os.Getenv(envLiteConnectionURL),
			TONStakingContractAddress: os.Getenv(envTONStakersContractAddr),
			KeyPairSecretName:         fmt.Sprintf("%s_%s", env, os.Getenv(envKeyPairSecretName)),
		}
	})

	return cfg
}
