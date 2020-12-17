package config

type Staking struct {
	DelegateVoteUnbondConfig DelegateVoteUnbondConfig `toml:"delegate_vote_unbond"`
	ProxyConfig              ProxyConfig              `toml:"proxy"`
}

type DelegateVoteUnbondConfig struct {
	DelegateNum string `toml:"delegate_num"`
	UnbondNum   string `toml:"unbond_num"`
	SleepTime   int    `toml:"sleep_time"`
}

type ProxyConfig struct {
	ProxyMnemonics []string `toml:"proxy_mnemonics"`
}
