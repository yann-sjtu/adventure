package market

import (
	"github.com/okex/adventure/x/strategy/order/market/cancel"
	"github.com/okex/adventure/x/strategy/order/market/maker"
	"github.com/okex/adventure/x/strategy/order/market/taker"
	"github.com/spf13/cobra"
)

func OrderMarketCmd() *cobra.Command {
	var orderMarketCmd = &cobra.Command{
		Use:   "order",
		Short: "cli about market service: maker & taker ",
	}
	orderMarketCmd.AddCommand(maker.MakerCmd())
	orderMarketCmd.AddCommand(taker.TakerCmd())
	orderMarketCmd.AddCommand(cancel.CancelCmd())
	return orderMarketCmd
}
