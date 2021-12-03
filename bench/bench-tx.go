package bench

import (
	"crypto/ecdsa"
	"log"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcmm "github.com/ethereum/go-ethereum/common"
	"github.com/okex/adventure/common"
	gosdk "github.com/okex/exchain-go-sdk"
	evmtypes "github.com/okex/exchain-go-sdk/module/evm/types"
	"github.com/spf13/cobra"
)

var (
	duplicate bool
	ethPort   int

	contract string
	direct   bool

	id    int64
	opts  []int64
	times int64
)

func OperateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operate",
		Short: "",
		Long:  "",
		Run:   startOperate,
	}
	flags := cmd.Flags()
	flags.IntVarP(&concurrency, "concurrency", "g", 10, "set the number of tx number per second")
	flags.IntVarP(&sleepTime, "sleepTime", "s", 1, "")

	flags.StringVarP(&privkPath, "privkey-path", "p", "", "")
	flags.StringSliceVarP(&rpc_hosts, "rpc-hosts", "u", []string{}, "")
	flags.StringVarP(&chainID, "chain-id", "i", "", "")

	flags.Int64Var(&id, "id", 0, "")
	flags.Int64SliceVar(&opts, "opts", []int64{}, "")
	flags.Int64Var(&times, "times", 0, "")

	flags.StringVar(&contract, "contract", "", "")
	flags.BoolVar(&direct, "direct", false, "")
	flags.BoolVar(&duplicate, "duplicate", false, "")
	flags.IntVar(&ethPort, "eth-port", 0, "if not zero, query on eth port 26659")
	return cmd
}

func startOperate(cmd *cobra.Command, args []string) {
	privkeys := common.GetPrivKeyFromPrivKeyFile(privkPath)
	param := Param{concurrency, rpc_hosts, chainID, privkeys, ethPort}

	RunTxs(param, operateFunc)
}

func operateFunc(cli *gosdk.Client, account *Account) {
	privateKey := account.GetPrivateKey()
	nonce := account.GetNonce()
	txdata := generateTxData()
	to := ethcmm.HexToAddress(contract)

	res, err := cli.Evm().SendTxEthereum(privateKey, nonce, to, nil, 2000000, evmtypes.DefaultGasPrice, txdata)
	if err != nil {
		log.Printf("send tx err: %s\n", err)
		if strings.Contains(err.Error(), "tx already exists in cache") {
			account.AddNonce()
		}
		time.Sleep(time.Second * time.Duration(sleepTime))
	} else {
		log.Printf("txhash: %s\n", res.TxHash)
		if duplicate {
			sendDuplicateTx(cli, privateKey, nonce, to, txdata)
		}
	}

	account.AddNonce()
	time.Sleep(time.Second * time.Duration(sleepTime))
}

func generateTxData() []byte {
	bigOpts := make([]*big.Int, len(opts))
	for i := 0; i < len(opts); i++ {
		bigOpts[i] = big.NewInt(opts[i])
	}

	if direct {
		operateABI, err := abi.JSON(strings.NewReader(OperateABI))
		if err != nil {
			panic(err)
		}
		txdata, err := operateABI.Pack("operate", bigOpts, big.NewInt(times))
		if err != nil {
			panic(err)
		}
		return txdata
	} else {
		routerABI, err := abi.JSON(strings.NewReader(RouterABI))
		if err != nil {
			panic(err)
		}
		txdata, err := routerABI.Pack("operate", big.NewInt(id), bigOpts, big.NewInt(times))
		if err != nil {
			panic(err)
		}
		return txdata
	}
}

func sendDuplicateTx(cli *gosdk.Client, privateKey *ecdsa.PrivateKey, nonce uint64, to ethcmm.Address, txdata []byte) {
	rand.Seed(time.Now().Unix())
	if rand.Intn(100) < 40 { // 40% chance to send duplicate txs
		for i := 1; i <= rand.Intn(3)+1; i++ {
			gasPrice := big.NewInt(1).Mul(evmtypes.DefaultGasPrice, big.NewInt(int64(i)))

			res, err := cli.Evm().SendTxEthereum(privateKey, nonce, to, nil, 2000000, gasPrice, txdata)
			if err != nil {
				log.Printf("[duplicate] error: %s\n", err)
			} else {
				log.Printf("[duplicate] txhash: %s\n", res.TxHash)
			}
		}
	}
}
