package multiwmt

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"io"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

type acc struct {
	privateKey string
	ecdsaPriv  *ecdsa.PrivateKey
	ethAddress common.Address
}

type wmtManager struct {
	clientList  []*ethclient.Client
	contracList []SwapContract
	superAcc    *acc

	worker  []*acc
	paraNum int
}

func newManager(cList []SwapContract, superAcc *acc, workPath string, paraNum int, clients []*ethclient.Client) *wmtManager {
	m := &wmtManager{
		clientList:  clients,
		contracList: cList,
		superAcc:    superAcc,
		paraNum:     paraNum,
	}
	m.prePareWorker(workPath)
	m.displayDetail()
	return m
}

func (m *wmtManager) displayDetail() {
	fmt.Println("contract size", len(m.contracList), "worker size:", len(m.worker), "paraNum", m.paraNum, "clientNum", len(m.clientList))
}

func (m *wmtManager) prePareWorker(path string) {
	f, err := os.Open(path)
	panicerr(err)
	defer f.Close()

	accList := make([]*acc, 0)
	rd := bufio.NewReader(f)
	for true {
		privKey, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}

		acc := keyToAcc(strings.TrimSpace(privKey))
		accList = append(accList, acc)
	}
	m.worker = accList
}

func GetNonce(client *ethclient.Client, privateKey *ecdsa.PrivateKey) uint64 {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("GetNonce Failed")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	cnt := 0
	for cnt < 50 {
		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			time.Sleep(2 * time.Second)
		} else {
			return nonce
		}
		cnt++
	}
	panic("GetNonce Failed")

}

func (m *wmtManager) Loop() {
	fmt.Printf("begin send wmt")

	var wg sync.WaitGroup
	workerIndex := 0

	for index := 0; index < m.paraNum; index++ {
		workIndexList := make([]int, 0)
		for i := 0; i < len(m.worker)/m.paraNum; i++ {
			workIndexList = append(workIndexList, workerIndex)
			workerIndex++
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			m.run(workIndexList)
		}()
	}
	wg.Wait()
}

func (m *wmtManager) needTransferToWorker() bool {
	return false
	// TODO : later to fix
	//for _, acc := range m.worker {
	//	bal, err := client.BalanceAt(context.Background(), acc.ethAddress, nil)
	//	if err == nil && bal.Int64() == 0 {
	//		return true
	//	}
	//}
	//return false
}

var (
	ether = new(big.Int).Mul(new(big.Int).SetInt64(1000000000), new(big.Int).SetInt64(1000000000))
)

func display(client *ethclient.Client, acc *acc, to common.Address, payload []byte) {
	data, err := client.CallContract(context.Background(), ethereum.CallMsg{
		From:     acc.ethAddress,
		To:       &to,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     payload,
	}, nil)
	if err == nil {
		fmt.Println("addr", acc.ethAddress.String(), "token balance", new(big.Int).SetBytes(data).String())
	} else {
		fmt.Println("err", err)
	}
}

func (m *wmtManager) DisPlayToken() {
	for _, acc := range m.worker {
		for _, c := range m.contracList {
			payload, err := erc20Builder.Build("balanceOf", acc.ethAddress)
			panicerr(err)
			display(m.clientList[0], acc, c.Token0, payload)
			display(m.clientList[0], acc, c.Token2, payload)
		}
	}
}

func (m *wmtManager) TransferToken0ToAccount() {
	needTransferToWorker := m.needTransferToWorker()

	nonce := GetNonce(m.clientList[0], m.superAcc.ecdsaPriv)
	fmt.Println("Begin TransferToken0ToAccount", "transferOkT:", needTransferToWorker)
	txs := make([]*types.Transaction, 0)
	for _, acc := range m.worker {
		for _, c := range m.contracList {
			if needTransferToWorker {
				tx := transferOkt(m.superAcc.privateKey, acc.ethAddress, nonce, ether)
				nonce++
				txs = append(txs, tx)
			}

			payload, err := erc20Builder.Build("transfer", acc.ethAddress, new(big.Int).SetInt64(10000000000))
			panicerr(err)
			tx := SignTxWithNonce(m.superAcc.ecdsaPriv, c.Token0, payload, nonce)
			nonce++
			txs = append(txs, tx)

			tx = SignTxWithNonce(m.superAcc.ecdsaPriv, c.Token2, payload, nonce)
			nonce++
			txs = append(txs, tx)
		}

	}
	fmt.Println("sendTx", len(txs), "use one node,may slow")
	if err := SendTxs(m.clientList[0], txs); err != nil {
		panic(err)
	}
	fmt.Println("end transferToken0ToAccount")
}

func (m *wmtManager) randomContract() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(5)
}

func (m *wmtManager) runPool(poolIndex int, workIndex int, contractIndex int, getReward bool) error {
	a := m.worker[workIndex]
	c := m.contracList[contractIndex]
	fmt.Println("run---", "workerIndex", workIndex, "contractIndex", contractIndex)

	token0 := c.Token0
	token1 := c.Token1
	lp := c.Lp1
	stakeRewards := c.StakingRewards1

	if poolIndex == 1 {
		token0 = c.Token2
		token1 = c.Token3
		lp = c.Lp2
		stakeRewards = c.StakingRewards2
	}

	nonce := GetNonce(m.clientList[workIndex%len(m.clientList)], a.ecdsaPriv)
	txList := make([]*types.Transaction, 0)

	// approve token0
	payload, err := erc20Builder.Build("approve", c.Router, new(big.Int).SetInt64(1000))
	panicerr(err)
	txList = append(txList, SignTxWithNonce(a.ecdsaPriv, token0, payload, nonce))
	nonce++

	// swap token0->token1
	payload, err = routerBuilder.Build("swapExactTokensForTokens",
		new(big.Int).SetInt64(500), new(big.Int),
		[]common.Address{token0, token1},
		a.ethAddress, big.NewInt(1956981781),
	)
	panicerr(err)
	txList = append(txList, SignTxWithNonce(a.ecdsaPriv, c.Router, payload, nonce))
	nonce++

	// approve token0 (for addLiquidity)
	payload, err = erc20Builder.Build("approve", c.Router, new(big.Int).SetInt64(30))
	panicerr(err)
	txList = append(txList, SignTxWithNonce(a.ecdsaPriv, token0, payload, nonce))
	nonce++

	// approve token1 (for addLiquidity)
	payload, err = erc20Builder.Build("approve", c.Router, new(big.Int).SetInt64(30))
	panicerr(err)
	txList = append(txList, SignTxWithNonce(a.ecdsaPriv, token1, payload, nonce))
	nonce++

	// addLiquidity
	payload, err = routerBuilder.Build("addLiquidity", token0, token1, new(big.Int).SetInt64(30), new(big.Int).SetInt64(30), new(big.Int), new(big.Int), a.ethAddress, big.NewInt(1956981781))
	panicerr(err)
	txList = append(txList, SignTxWithNonce(a.ecdsaPriv, c.Router, payload, nonce))
	nonce++

	// approve lp for stakingRewards
	payload, err = erc20Builder.Build("approve", stakeRewards, new(big.Int).SetInt64(10))
	panicerr(err)
	txList = append(txList, SignTxWithNonce(a.ecdsaPriv, lp, payload, nonce))
	nonce++

	// stake lp for stakingRewards
	payload, err = StakingRewardsBuilder.Build("stake", big.NewInt(10))
	panicerr(err)
	txList = append(txList, SignTxWithNonce(a.ecdsaPriv, stakeRewards, payload, nonce))
	nonce++

	if getReward {
		// getReward for stakingRewards
		payload, err = StakingRewardsBuilder.Build("getReward")
		panicerr(err)
		txList = append(txList, SignTxWithNonce(a.ecdsaPriv, stakeRewards, payload, nonce))
		nonce++
	}

	payload, err = StakingRewardsBuilder.Build("withdraw", big.NewInt(3))
	panicerr(err)
	txList = append(txList, SignTxWithNonce(a.ecdsaPriv, stakeRewards, payload, nonce))
	nonce++

	if getReward {
		payload, err = StakingRewardsBuilder.Build("exit")
		panicerr(err)
		txList = append(txList, SignTxWithNonce(a.ecdsaPriv, stakeRewards, payload, nonce))
		nonce++
	}

	if err := SendTxs(m.clientList[workIndex%len(m.clientList)], txList); err != nil {
		return err
	}

	time.Sleep(2 * time.Second)
	return nil
}

func (m *wmtManager) run(tasks []int) {
	rand.Seed(time.Now().UnixNano())
	sleepTime := rand.Intn(10)
	time.Sleep(time.Duration(sleepTime) * time.Second)
	turns := 0
	for true {
		for _, workIndex := range tasks {
			getReward := turns%2 == 1
			for contractIndex, _ := range m.contracList {
				if err := m.runPool(0, workIndex, contractIndex, getReward); err != nil {
					fmt.Println("runErr-0", workIndex)
					continue
				}

				if err := m.runPool(1, workIndex, contractIndex, getReward); err != nil {
					fmt.Println("runErr-1", workIndex)
					continue
				}
			}
		}
		turns++
	}
}
