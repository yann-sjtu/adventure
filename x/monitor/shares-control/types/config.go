package types

var Cfg *Config

type Config struct {
	targetValAddrs   []string `toml:"target_validator_addresses"`
	workersAccInfo   []string `toml:"worker_infos"`
	percentToPlunder string   `toml:"rewards_percentage"`
}
