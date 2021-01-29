package types

var Cfg *Config

type Config struct {
	TargetValAddrs    []string `toml:"target_validator_addresses"`
	WorkersAccInfo    []string `toml:"workers_infos"`
	PercentToDominate string   `toml:"domination_percentage"`
}
