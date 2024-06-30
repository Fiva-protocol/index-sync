package main

import (
	"flag"
	"sync"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type (
	Config struct {
		sync.RWMutex
		MasterContractAddress     string `mapstructure:"master_contract_address"`
		LiteConnectionsURLTestnet string `mapstructure:"lite_connections_url_testnet"`
		LiteConnectionsURLMainnet string `mapstructure:"lite_connections_url_mainnet"`
		TONStakingContractAddress string `mapstructure:"ton_stakers_contract_address"`

		AWSConfig `mapstructure:",squash"`
	}

	AWSConfig struct {
		SecretManagerSecretName string `mapstructure:"aws_secret_manager_secret_name"`
		AccessKeyID             string `mapstructure:"aws_access_key_id"`
		SecretAccessKey         string `mapstructure:"aws_secret_access_key"`
		Region                  string `mapstructure:"aws_region"`
	}
)

var envs = []*EnvVar{
	DefaultMnemonic,
	DefaultMasterContractAddress,
	DefaultLiteConnectionsURLTestnet,
	DefaultLiteConnectionsURLMainnet,
	DefaultTonStakingContractAddress,

	// AWS
	DefaultAWSSecretName,
	DefaultAWSAccessKeyID,
	DefaultAWSSecretAccessKey,
	DefaultAWSRegion,
}

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
