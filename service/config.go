package main

import (
	"flag"
	"sync"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultLiteConnectionURL         = "https://ton-blockchain.github.io/testnet-global.config.json"
	defaultTONStakersContractAddress = "EQD2_4d91M4TVbEBVyBF8J1UwpMJc361LKVCz6bBlffMW05o"
)

type (
	Config struct {
		sync.RWMutex
		LiteConnectionsURL        string `mapstructure:"lite_connections_url_testnet"`
		TONStakingContractAddress string `mapstructure:"ton_stakers_contract_address"`

		AWSConfig `mapstructure:",squash"`
	}

	AWSConfig struct {
		SecretManagerSecretName string `mapstructure:"aws_secret_manager_secret_name"`
		AccessKeyID             string `mapstructure:"aws_access_key_id"`
		SecretAccessKey         string `mapstructure:"aws_secret_access_key"`
		Region                  string `mapstructure:"aws_region"`
	}

	EnvVar struct {
		DefaultValue           interface{}
		Flag, Env, Description string
	}
)

var (
	envs = []*EnvVar{
		DefaultMnemonic,
		DefaultLiteConnectionsURL,
		DefaultTonStakingContractAddress,
	}

	DefaultMnemonic = NewEnvVar(
		"mnemonic",
		"MNEMONIC",
		"",
		"Seed phrase",
	)

	DefaultTonStakingContractAddress = NewEnvVar(
		"ton_stakers_contract_address",
		"TON_STAKERS_CONTRACT_ADDRESS",
		defaultTONStakersContractAddress,
		"TON stakers contract address",
	)

	DefaultLiteConnectionsURL = NewEnvVar(
		"lite_connections_url",
		"LITE_CONNECTIONS_URL",
		defaultLiteConnectionURL,
		"URL to lite connections",
	)
)

func NewConfig() (*Config, error) {
	c := Config{}
	if err := GetConfig(&c, []*EnvVar{}); err != nil {
		return nil, err
	}

	return &c, nil
}

func BindConfig() {
	for _, e := range envs {
		switch val := e.DefaultValue.(type) {
		case string:
			flag.String(e.Flag, val, e.Description)
		case int:
			flag.Int(e.Flag, val, e.Description)
		case bool:
			flag.Bool(e.Flag, val, e.Description)
		case uint64:
			flag.Uint64(e.Flag, val, e.Description)
		case time.Duration:
			flag.Duration(e.Flag, val, e.Description)
		default:
			continue
		}
		if e.DefaultValue != nil {
			viper.SetDefault(e.Env, e.DefaultValue)
		}
	}

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}

	var err error
	for _, e := range envs {
		err = viper.BindEnv(e.Env)
		if err != nil {
			panic(err)
		}
	}
}

func AddEnvs(customEnvs []*EnvVar) {
	var tmpEnvs []*EnvVar
	tmpEnvs = append(tmpEnvs, customEnvs...)
	for _, defaultEnv := range envs {
		check := true
		for _, customEnv := range customEnvs {
			if customEnv.Flag == defaultEnv.Flag {
				check = false
				break
			}
		}

		if check {
			tmpEnvs = append(tmpEnvs, defaultEnv)
		}
	}

	envs = tmpEnvs
}

func GetConfig(cfg interface{}, customEnvs []*EnvVar) error {
	AddEnvs(customEnvs)
	BindConfig()
	return viper.Unmarshal(cfg)
}

func NewEnvVar(flag, env string, defaultValue interface{}, description string) *EnvVar {
	return &EnvVar{Flag: flag, Env: env, Description: description, DefaultValue: defaultValue}
}
