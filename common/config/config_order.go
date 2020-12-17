package config

type Order struct {
	NewConfig NewConfig `toml:"new"`
}

type NewConfig struct {
	Products     []string `toml:"products"`
	BuyPrice     string   `toml:"buy_price"`
	SellPrice    string   `toml:"sell_price"`
	BuyQuantity  string   `toml:"buy_quantity"`
	SellQuantity string   `toml:"sell_quantity"`
}
