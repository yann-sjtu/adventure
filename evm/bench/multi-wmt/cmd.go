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
		Short: "evm cli of test strategy",
		Args:  cobra.NoArgs,
		Run:   wmt,
	}
	return wmtCmd
}

var (
	chainID      = new(big.Int).SetUint64(65)
	signer       = types.NewEIP155Signer(chainID)
	gasPrice     = new(big.Int).SetUint64(1000000000)
	gasLimit     = uint64(3000000)
	useOldTxHash = bool(false)
)

func wmt(cmd *cobra.Command, args []string) {
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
	m := newManager(cList, superAcc, c.WorkerPath, c.ParaNum, clients)

	m.TransferToken0ToAccount()
	m.Loop()
}
