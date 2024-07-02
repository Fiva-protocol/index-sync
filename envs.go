package main

const (
	defaultAddress = "EQAHmEAgPST8XV4GN6r6E4NesuLs7lDbzsSW1ougMxItut9S"
	// defaultAddress                   = "EQCnF65lKoXuXxq4XWDVT8OHA6C4XUsrGYjgkCsBpBJKgL51"
	defaultLiteConnectionURLTestnet  = "https://ton-blockchain.github.io/testnet-global.config.json"
	defaultLiteConnectionURLMainnet  = "https://ton.org/global.config.json"
	defaultTONStakersContractAddress = "EQD2_4d91M4TVbEBVyBF8J1UwpMJc361LKVCz6bBlffMW05o"
	defaultAWSRegion                 = "us-east-1"
)

type EnvVar struct {
	DefaultValue           interface{}
	Flag, Env, Description string
}

var (
	DefaultMnemonic = NewEnvVar(
		"mnemonic",
		"MNEMONIC",
		"",
		"Seed phrase",
	)

	DefaultMasterContractAddress = NewEnvVar(
		"master_contract_address",
		"MASTER_CONTRACT_ADDRESS",
		defaultAddress,
		"Master contract address",
	)

	DefaultTonStakingContractAddress = NewEnvVar(
		"ton_stakers_contract_address",
		"TON_STAKERS_CONTRACT_ADDRESS",
		defaultTONStakersContractAddress,
		"TON stakers contract address",
	)

	DefaultLiteConnectionsURLTestnet = NewEnvVar(
		"lite_connections_url_testnet",
		"LITE_CONNECTIONS_URL_TESTNET",
		defaultLiteConnectionURLTestnet,
		"URL to lite connections testnet",
	)

	DefaultLiteConnectionsURLMainnet = NewEnvVar(
		"lite_connections_url_mainnet",
		"LITE_CONNECTIONS_URL_MAINNET",
		defaultLiteConnectionURLMainnet,
		"URL to lite connections mainnet",
	)

	DefaultAWSSecretName = NewEnvVar(
		"aws_secret_manager_secret_name",
		"AWS_SECRET_MANAGER_SECRET_NAME",
		"",
		"AWS Secret Manager secret name",
	)

	DefaultAWSAccessKeyID = NewEnvVar(
		"aws_access_key_id",
		"AWS_ACCESS_KEY_ID",
		"",
		"AWS Access Key ID",
	)

	DefaultAWSSecretAccessKey = NewEnvVar(
		"aws_secret_access_key",
		"AWS_SECRET_ACCESS_KEY",
		"",
		"AWS Secret Manager secret name",
	)

	DefaultAWSRegion = NewEnvVar(
		"aws_region",
		"AWS_REGION",
		defaultAWSRegion,
		"AWS Secret Manager secret name",
	)
)

func NewEnvVar(flag, env string, defaultValue interface{}, description string) *EnvVar {
	return &EnvVar{Flag: flag, Env: env, Description: description, DefaultValue: defaultValue}
}
