package main

const (
	defaultAddress           = "EQCnF65lKoXuXxq4XWDVT8OHA6C4XUsrGYjgkCsBpBJKgL51"
	defaultLiteConnectionURL = "https://ton-blockchain.github.io/testnet-global.config.json"
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
		"mnemonic_contract_address",
		"MNEMONIC_CONTRACT_ADDRESS",
		defaultAddress,
		"Host. Related env var",
	)

	DefaultLiteConnectionsURL = NewEnvVar(
		"lite_connections_url",
		"LITE_CONNECTIONS_URL",
		defaultLiteConnectionURL,
		"Host. Related env var",
	)
)

func NewEnvVar(flag, env string, defaultValue interface{}, description string) *EnvVar {
	return &EnvVar{Flag: flag, Env: env, Description: description, DefaultValue: defaultValue}
}
