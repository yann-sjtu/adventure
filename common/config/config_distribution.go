package config

type Distribution struct {
	SetWithdrawAddrConfig SetWithdrawAddrConfig `toml:"set_withdraw_addr"`
}

type SetWithdrawAddrConfig struct {
	Address []string `toml:"address"`
}
