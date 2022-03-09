package client

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthClient struct {
	*ethclient.Client
	signer types.Signer
}

func NewEthClient(ip string) (*EthClient, error) {
	cli, err := ethclient.Dial(ip)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize client: %+v", err)
	}

	chainId, err := cli.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	return &EthClient{
		cli,
		types.NewLondonSigner(chainId),
	}, nil
}

func (e EthClient) QueryNonce(hexAddr string) (uint64, error) {
	nonce, err := e.PendingNonceAt(context.Background(), ethcmn.HexToAddress(hexAddr))
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

func (e EthClient) SendEthereumTx(privatekey *ecdsa.PrivateKey, nonce uint64, to ethcmn.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) (ethcmn.Hash, error) {
	// 1. make tx
	unsignedTx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data)

	// 2. sign unsignedTx -> rawTx
	signedTx, err := types.SignTx(unsignedTx, e.signer, privatekey)
	if err != nil {
		return ethcmn.Hash{}, err
	}

	// 3. send rawTx
	err = e.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return ethcmn.Hash{}, err
	}

	return signedTx.Hash(), err
}

func (e EthClient) CreateContract(privatekey *ecdsa.PrivateKey, nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) (ethcmn.Hash, error) {
	// 1. make tx
	unsignedTx := types.NewContractCreation(nonce, amount, gasLimit, gasPrice, data)

	// 2. sign unsignedTx -> rawTx
	signedTx, err := types.SignTx(unsignedTx, e.signer, privatekey)
	if err != nil {
		return ethcmn.Hash{}, err
	}

	// 3. send rawTx
	err = e.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return ethcmn.Hash{}, err
	}

	return signedTx.Hash(), err
}

func (e EthClient) GetEthTxRlpEncode(pk *ecdsa.PrivateKey, nonce uint64, to ethcmn.Address, amount *big.Int, gaslimit uint64, gasprice *big.Int, data []byte)(string, error){
	//make tx
	unsignedTx := types.NewTransaction(nonce,to,amount,gaslimit,gasprice,data)

	//sign tx
	signedTx, err := types.SignTx(unsignedTx, e.signer, pk)
	if err != nil {
		return ethcmn.Hash{}.String(), err
	}

	//当需要调用 eth_sendRawTransaction 函数中的 params的时候，通过下面这个rlp来构造
	b, err := rlp.EncodeToBytes(signedTx)
	params := "0x" + hex.EncodeToString(b)

	return params, nil

}


