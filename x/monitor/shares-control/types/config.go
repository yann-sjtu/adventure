package types

var Cfg *Config

type Config struct {
	targetValAddrs   []string
	workersAccInfo   []string
	valNumberInTop21 string
	percentToPlunder string
}
