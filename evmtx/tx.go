package evmtx

import (
	"crypto/ecdsa"
	"log"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	routerAddr = "0x2CA0E1278B9D7A967967d3C707b81C72FC180CaF"
	wethAddr   = "0x70c1c53E991F31981d592C2d865383AC0d212225"
	usdtAddr   = "0xee666e967293094007d7c3718044e07565b1f8a9"
	dyfAddr    = "0xd78e1680e93bF57580F299d75B364e591873a8d3"
)

var (
	defaultGasPrice, _ = sdk.ParseDecCoin("0.000000001okt")

	depositPayloadStr string
	approvePayloadStr string
	swapPayloadStr    string
	dyfPayloadStr     string
)

func getParam(txtype int) map[string]string {
	param := make(map[string]string)
	switch txtype {
	case 0: // deposit
		param["to"] = wethAddr
		param["value"] = "0x2386f26fc10000"
		param["data"] = depositPayloadStr
	case 1: //approve
		param["to"] = wethAddr
		param["value"] = "0x0"
		param["data"] = approvePayloadStr
	case 2: //swap
		param["to"] = routerAddr
		param["value"] = "0x0"
		param["data"] = swapPayloadStr
	case 3: //dyf
		param["to"] = dyfAddr
		param["value"] = "0x0"
		param["data"] = dyfPayloadStr
	}

	return param
}

func generateTxParams(from string, txtype int) []map[string]string {
	param := make([]map[string]string, 1)
	param[0] = getParam(txtype)
	param[0]["from"] = from
	param[0]["gasPrice"] = (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String()

	return param
}

func signTx(privKey *ecdsa.PrivateKey, nonce hexutil.Uint64, to string, amount string, gas hexutil.Uint64, gasPrice hexutil.Big, data string) string {
	unsignedTx := ethtypes.NewTransaction(
		uint64(nonce),
		ethcmn.HexToAddress(to),
		hexutil.MustDecodeBig(amount),
		uint64(gas),
		(*big.Int)(&gasPrice),
		hexutil.MustDecode(data),
	)

	signedTx, err := ethtypes.SignTx(unsignedTx, ethtypes.NewEIP155Signer(big.NewInt(65)), privKey)
	if err != nil {
		log.Fatalf("failed to sign the unsignedTx offline: %+v", err)
	}
	rawdata, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		panic(err)
	}
	return hexutil.Encode(rawdata)
}
