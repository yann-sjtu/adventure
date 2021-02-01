package types

var Cfg *Config

type Config struct {
	OurValAddrs      []string `toml:"our_validator_addresses"`
	OurTop18Addrs    []string `toml:"our_top_18_addresses"`
	WorkersAccInfo   []string `toml:"workers_infos"`
	PercentToPlunder string   `toml:"plundered_percentage"`
}
