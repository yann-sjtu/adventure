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
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2"
)

var (
	defaultGasPrice, _ = sdk.ParseDecCoin("0.000000001okt")
)

func generateTxParams(from string) []map[string]string {
	param := make([]map[string]string, 1)
	param[0] = map[string]string{
		"from":     from,
		"to":       "0xEe666E967293094007d7C3718044e07565B1f8a9", // todo
		"value":    "0x0",
		"gasPrice": (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String(),
		"data":     "0x" + ethcmn.Bytes2Hex(UniswapV2.BuildApprovePayload("0x1bF1D0fF0dBE2A3ffD4E287033E3354e9Ad5Ea80", 100)), //todo
	}

	return param
}

func signTx(privKey *ecdsa.PrivateKey, nonce uint64, gas uint64, gasPrice *big.Int) string {
	unsignedTx := ethtypes.NewTransaction(
		nonce,
		ethcmn.HexToAddress("0xEe666E967293094007d7C3718044e07565B1f8a9"),
		big.NewInt(0),
		gas,
		gasPrice,
		UniswapV2.BuildApprovePayload("0x1bF1D0fF0dBE2A3ffD4E287033E3354e9Ad5Ea80", 100),
	)

	signedTx, err := ethtypes.SignTx(unsignedTx, ethtypes.NewEIP155Signer(big.NewInt(65)), privKey)
	if err != nil {
		log.Fatalf("failed to sign the unsignedTx offline: %+v", err)
	}
	data, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		panic(err)
	}
	return hexutil.Encode(data)
}
