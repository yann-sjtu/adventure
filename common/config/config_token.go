package config

type Token struct {
	MultiSendConfig MultiSendConfig `toml:"multi_send"`
}

type MultiSendConfig struct {
	ToAddrsPath string `toml:"to_addrs_path"`
}
