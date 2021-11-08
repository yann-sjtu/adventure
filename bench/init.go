package bench

import (
	"log"
	"time"

	"github.com/okex/adventure/common"
	gosdk "github.com/okex/exchain-go-sdk"
	"github.com/okex/exchain-go-sdk/types"
	"github.com/spf13/cobra"
)

var (
	containerContract string
	//addresses []ethcmm.Address
)

func InitStorageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init-storage",
		Short: "",
		Long:  "",
		Run:   initStorage,
	}
	flags := cmd.Flags()
	flags.IntVarP(&concurrency, "concurrency","g", 10, "set the number of tx number per second")
	flags.IntVarP(&sleepTimeTx, "sleepTime", "s",1, "")

	flags.StringVarP(&privkPath, "privkey-path", "p", "","")
	flags.StringSliceVarP(&rpc_hosts, "rpc-hosts","u", []string{}, "")
	flags.StringVarP(&chainID, "chain-id", "i","", "")

	flags.StringVar(&containerContract, "container-contract","0xa1ddCC79DAAf7d3bE05E83f8d583EE353713cAe0", "")
	return cmd
}

func initStorage(cmd *cobra.Command, args []string) {
	//addresses = readAddress()

	privkeys := common.GetPrivKeyFromPrivKeyFile(privkPath)
	for i := 0; i < concurrency; i++ {
		go func(index int, privkey string) {
			rpcHost := rpc_hosts[index%len(rpc_hosts)]
			deploy(privkey, rpcHost)
		}(i, privkeys[i])
	}

	select {}
}

func deploy(privkey string, host string) {
	cfg, _ := types.NewClientConfig(host, chainID, types.BroadcastSync, "", 30000000, 1.5, "0.0000000001"+common.NativeToken)
	cli := gosdk.NewClient(cfg)

	addr := getCosmosAddress(privkey)
	accInfo, err := cli.Auth().QueryAccount(addr.String())
	if err != nil {
		panic(err)
	}
	nonce := accInfo.GetSequence()

	for {
		res, err := cli.Evm().SendTxEthereum(privkey, containerContract, "", "0xfe4b84df0000000000000000000000000000000000000000000000000000000000000002",20000000, nonce)
		if err != nil {
			continue
		} else {
			log.Printf("txhash: %s\n", res.TxHash)
		}

		nonce++
		time.Sleep(time.Second * time.Duration(sleepTimeTx))
	}
}
//
//func appends(privkey string, host string) {
//	cfg, _ := types.NewClientConfig(host, chainID, types.BroadcastSync, "", 30000000, 1.5, "0.0000000001"+common.NativeToken)
//	cli := gosdk.NewClient(cfg)
//
//	addr := getCosmosAddress(privkey)
//	accInfo, err := cli.Auth().QueryAccount(addr.String())
//	if err != nil {
//		panic(err)
//	}
//	nonce := accInfo.GetSequence()
//
//	containerABI, err := abi.JSON(strings.NewReader(ContainerJson))
//	if err != nil {
//		panic(err)
//	}
//
//	for k := 0; k < 100900; k++ {
//		if k % 100 == 0 {
//			addrs := addresses[k:k+100]
//			txdata, err := containerABI.Pack("append", addrs)
//			if err != nil {
//				panic(err)
//			}
//
//			res, err := cli.Evm().SendTxEthereum(privkey, containerContract, "", ethcmm.Bytes2Hex(txdata),25000000, nonce)
//			if err != nil {
//				continue
//			} else {
//				log.Printf("txhash: %s\n", res.TxHash)
//			}
//
//			nonce++
//			time.Sleep(time.Second * time.Duration(sleepTimeTx))
//		}
//	}
//}
//
//func readAddress() []ethcmm.Address {
//	f, err := os.Open("/Users/green/project/okex/ethwsclient/address")
//	if err != nil {
//		log.Fatalln(err.Error())
//		return nil
//	}
//	defer f.Close()
//
//	addrs := make([]ethcmm.Address , 100904, 100904)
//	rd := bufio.NewReader(f)
//	for index := 0; index < 100905; index++ {
//		privKey, err := rd.ReadString('\n')
//		if err != nil || io.EOF == err {
//			break
//		}
//		addrs[index] = ethcmm.HexToAddress(strings.TrimSpace(privKey))
//	}
//
//	return addrs
//}
//
//var ContainerJson = `
//[
//	{
//		"inputs": [
//			{
//				"internalType": "address[]",
//				"name": "_storages",
//				"type": "address[]"
//			}
//		],
//		"name": "append",
//		"outputs": [],
//		"stateMutability": "nonpayable",
//		"type": "function"
//	},
//	{
//		"inputs": [
//			{
//				"internalType": "uint256",
//				"name": "num",
//				"type": "uint256"
//			}
//		],
//		"name": "initialize",
//		"outputs": [],
//		"stateMutability": "nonpayable",
//		"type": "function"
//	},
//	{
//		"inputs": [],
//		"name": "length",
//		"outputs": [
//			{
//				"internalType": "uint256",
//				"name": "",
//				"type": "uint256"
//			}
//		],
//		"stateMutability": "view",
//		"type": "function"
//	},
//	{
//		"inputs": [
//			{
//				"internalType": "uint256",
//				"name": "",
//				"type": "uint256"
//			}
//		],
//		"name": "storages",
//		"outputs": [
//			{
//				"internalType": "address",
//				"name": "",
//				"type": "address"
//			}
//		],
//		"stateMutability": "view",
//		"type": "function"
//	}
//]`