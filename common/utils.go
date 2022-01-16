package common

import (
	"crypto/ecdsa"
	"fmt"

	ethcmm "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func GetEthAddressFromPK(privateKey *ecdsa.PrivateKey) ethcmm.Address {
	pubkeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if ok != true {
		panic(fmt.Errorf("convert into pubkey failed"))
	}
	return crypto.PubkeyToAddress(*pubkeyECDSA)
}
