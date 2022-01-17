package multiwmt

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"io"
	"math/big"
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

type group struct {
	contract SwapContract
	worker   []*acc
}

type wmtManager struct {
	contracList []SwapContract
	superAcc    *acc

	worker          []*acc
	workMapContract map[int]int
	paraNum         int
}

func newManager(cList []SwapContract, superAcc *acc, workPath string, paraNum int) *wmtManager {
	m := &wmtManager{
		contracList:     cList,
		superAcc:        superAcc,
		workMapContract: make(map[int]int, 0),
		paraNum:         paraNum,
	}
	m.prePareWorker(workPath)
	m.calGroup()
	m.displayDetail()
	return m
}

func (m *wmtManager) displayDetail() {
	fmt.Println("contract size", len(m.contracList), "worker size:", len(m.worker), "paraNum", m.paraNum)
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

func (m *wmtManager) calGroup() {
	contractIndex := 0
	for index, _ := range m.worker {
		m.workMapContract[index] = contractIndex % len(m.contracList)
		contractIndex++
	}
}

func GetNonce(privateKey *ecdsa.PrivateKey) uint64 {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("GetNonce Failed")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	cnt := 0
	for cnt < 10 {
		time.Sleep(3 * time.Second)
		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {

		} else {
			return nonce
		}
		cnt++
	}
	panic("GetNonce Failed")

}

func (m *wmtManager) Loop() {

	h, err := client.HeaderByNumber(context.Background(), nil)
	panicerr(err)
	fmt.Printf("begin send wmt : currBlockHeight:%d goRountineNums:%d\n", h.Number, m.paraNum)

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
	for _, acc := range m.worker {
		bal, err := client.BalanceAt(context.Background(), acc.ethAddress, nil)
		if err == nil && bal.Int64() == 0 {
			return true
		}
	}
	return false
}

var (
	ether = new(big.Int).Mul(new(big.Int).SetInt64(1000000000), new(big.Int).SetInt64(1000000000))
)

func (m *wmtManager) TransferToken0ToAccount() {
	needTransferToWorker := m.needTransferToWorker()

	nonce := GetNonce(m.superAcc.ecdsaPriv)
	fmt.Println("Begin TransferToken0ToAccount", "transferOkT:", needTransferToWorker)
	txs := make([]*types.Transaction, 0)
	accCnt := 0
	for index, acc := range m.worker {
		c := m.contracList[m.workMapContract[index]]
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
		accCnt++

	}
	fmt.Println("SendTx")
	if err := SendTxs(txs); err != nil {
		panic(err)
	}
	fmt.Println("GetReceipt")
	if err := getReceipt(txs); err != nil {
		panic(err)
	}
	fmt.Println("End TransferToken0ToAccount")
}

func (m *wmtManager) runPool(poolIndex int, c SwapContract, a *acc) error {
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

	nonce := GetNonce(a.ecdsaPriv)
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

	// getReward for stakingRewards
	payload, err = StakingRewardsBuilder.Build("getReward")
	panicerr(err)
	txList = append(txList, SignTxWithNonce(a.ecdsaPriv, stakeRewards, payload, nonce))
	nonce++

	payload, err = StakingRewardsBuilder.Build("withdraw", big.NewInt(3))
	panicerr(err)
	txList = append(txList, SignTxWithNonce(a.ecdsaPriv, stakeRewards, payload, nonce))
	nonce++

	payload, err = StakingRewardsBuilder.Build("exit")
	panicerr(err)
	txList = append(txList, SignTxWithNonce(a.ecdsaPriv, stakeRewards, payload, nonce))
	nonce++
	if err := SendTxs(txList); err != nil {
		return err
	}
	return getReceipt(txList)
}

var (
	mapp = make(map[common.Address]bool)
)

func (m *wmtManager) run(tasks []int) {
	for true {
		for _, workIndex := range tasks {
			a := m.worker[workIndex]
			c := m.contracList[m.workMapContract[workIndex]]
			fmt.Println("run---", "contractIndex", m.workMapContract[workIndex], "workerIndex", workIndex)
			if err := m.runPool(0, c, a); err != nil {
				fmt.Println("runErr-0", c.Token0, a.ethAddress, err)
				continue
			}

			if err := m.runPool(1, c, a); err != nil {
				fmt.Println("runErr-1", c.Token0, a.ethAddress, err)
				continue
			}

		}
	}

}
