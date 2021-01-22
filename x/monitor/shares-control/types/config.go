package types

var Cfg *Config

type Config struct {
	TargetValAddrs    []string `toml:"target_validator_addresses"`
	WorkersAccInfo    []string `toml:"worker_infos"`
	PercentToPlunder  string   `toml:"rewards_percentage"`
	PercentToDominate string   `toml:"domination_percentage"`
}
