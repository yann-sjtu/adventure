package cancel

import (
	"strings"

	"github.com/okex/adventure/common"
	"github.com/okex/adventure/common/config"
	"github.com/okex/adventure/x/order/market/types"
	"github.com/okex/okexchain-go-sdk/utils"
	"github.com/spf13/cobra"
)

func CancelCmd() *cobra.Command {
	makerCmd := &cobra.Command{
		Use: "cancel",
		Run: runCancel,
	}
	flags := makerCmd.Flags()
	flags.StringVarP(&types.Product, "product", "p", "", "set coin product name")
	flags.StringVarP(&types.Mnemonic, "mnemonic", "m", "", "set account mnemonic")
	return makerCmd
}

func runCancel(cmd *cobra.Command, args []string) {
	// init maker account info
	info, _, err := utils.CreateAccountWithMnemo(types.Mnemonic, types.Name, types.PassWd)
	if err != nil {
		panic(err)
	}

	msg := types.QueryOrders(types.Product, info.GetAddress().String())
	var ids string
	i := 0
	for _, order := range msg.Data.Data {
		if i >= 200 {
			break
		}
		ids += order.OrderID + ","
		i++
	}

	// init clients
	clientManager := common.NewClientManager(config.Cfg.Hosts, config.AUTO)
	types.CancelOrders(clientManager.GetRandomClient(), info, strings.TrimRight(ids, ","))
}
