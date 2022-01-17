package multiwmt

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"math/big"
)

func MultiWmtCmt() *cobra.Command {
	var wmtCmd = &cobra.Command{
		Use:   "multiwmt",
		Short: "wmt-run",
		Args:  cobra.NoArgs,
		Run:   wmtRun,
	}
	return wmtCmd
}

func MultiWmtInit() *cobra.Command {
	var wmtCmd = &cobra.Command{
		Use:   "multiinit",
		Short: "wmt-init",
		Args:  cobra.NoArgs,
		Run:   wmtInit,
	}
	return wmtCmd
}

func MultiTokenBalance() *cobra.Command {
	var wmtCmd = &cobra.Command{
		Use:   "multitoken",
		Short: "wmt-token",
		Args:  cobra.NoArgs,
		Run:   wmtToken,
	}
	return wmtCmd
}

var (
	chainID  = new(big.Int).SetUint64(65)
	signer   = types.NewEIP155Signer(chainID)
	gasPrice = new(big.Int).SetUint64(1000000000)
	gasLimit = uint64(3000000)
)

func getM() *wmtManager {
	c := loadWMTConfig("./config/wmt.json")

	initBuilder()
	initClient(c)
	cList := LoadContractList(c.ContractPath)
	clients := make([]*ethclient.Client, 0)
	for _, v := range c.RPC {
		c, err := ethclient.Dial(v)
		panicerr(err)
		clients = append(clients, c)
	}
	superAcc := keyToAcc(c.SuperAcc)
	return newManager(cList, superAcc, c.WorkerPath, c.ParaNum, clients)
}
func wmtRun(cmd *cobra.Command, args []string) {
	m := getM()
	m.Loop()
}

func wmtInit(cmd *cobra.Command, args []string) {
	m := getM()
	m.TransferToken0ToAccount()
}

func wmtToken(cmd *cobra.Command, args []string) {
	m := getM()
	m.DisPlayToken()
}
