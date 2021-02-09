package evm

import "github.com/spf13/cobra"

func uniswapTestCmd() *cobra.Command {
	InitTemplate()

	cmd := &cobra.Command{
		Use:   "uniswap-testnet-operate",
		Short: "",
		Args:  cobra.NoArgs,
		Run:   testLoop,
	}

	//flags := cmd.Flags()
	//flags.IntVarP(&Num, "num", "n", 1000, "set Num of issusing token")
	//flags.IntVarP(&GoroutineNum, "goroutine-num", "g", 1, "set Goroutine Num of deploying contracts")
	//flags.IntVarP(&TransferGoNum, "transfer-go-num", "t", 1, "set Goroutine Num of transfering erc20 token")
	//flags.StringVarP(&MnemonicPath, "mnemonic-path", "m", "", "set the MnemonicPath path")

	return cmd
}

const (
	usdtAddr = "0xffea71957a3101d14474a3c358ede310e17c2409"
)

func testLoop(cmd *cobra.Command, args []string) {

}


