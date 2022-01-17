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

type group struct {
	contract SwapContract
	worker   []*acc
}

type wmtManager struct {
	contracList []SwapContract
	superAcc    *acc

	worker []*acc

	groupList []*group
	paraNum   int
}

func newManager(cList []SwapContract, superAcc *acc, workPath string, paraNum int) *wmtManager {
	m := &wmtManager{
		contracList: cList,
		superAcc:    superAcc,
		groupList:   make([]*group, 0),
		paraNum:     paraNum,
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
	workSizePreContract := len(m.worker) / len(m.contracList)
	currentWorkerIndex := 0
	for index := 0; index < len(m.contracList); index++ {
		tm := make([]*acc, 0)
		for i := 0; i < workSizePreContract; i++ {
			tm = append(tm, m.worker[currentWorkerIndex])
			currentWorkerIndex++
		}
		groupInstance := &group{
			worker:   tm,
			contract: m.contracList[index],
		}
		m.groupList = append(m.groupList, groupInstance)
	}
}

func GetNonce(key string) uint64 {
	privateKey := getPrivateKey(key)
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
	cList := make([]SwapContract, 0)
	accList := make([]*acc, 0)

	for _, group := range m.groupList {
		for _, acc := range group.worker {
			if len(cList) >= m.paraNum {
				break
			}
			cList = append(cList, group.contract)
			accList = append(accList, acc)
		}
	}

	h, err := client.HeaderByNumber(context.Background(), nil)
	panicerr(err)
	fmt.Printf("begin send wmt : currBlockHeight:%d goRountineNums:%d\n", h.Number, m.paraNum)

	var wg sync.WaitGroup
	for k, contract := range cList {
		wg.Add(1)
		cc := contract
		kk := k
		go func() {
			defer wg.Done()
			m.run(cc, accList[kk])
		}()
	}
	wg.Wait()
}

func (m *wmtManager) needTransferToWorker() bool {
	for _, group := range m.groupList {
		for _, acc := range group.worker {
			bal, err := client.BalanceAt(context.Background(), acc.ethAddress, nil)
			if err == nil && bal.Int64() == 0 {
				return true
			}

		}
	}
	return false
}

var (
	ether = new(big.Int).Mul(new(big.Int).SetInt64(1000000000), new(big.Int).SetInt64(1000000000))
)

func (m *wmtManager) TransferToken0ToAccount() {
	needTransferToWorker := m.needTransferToWorker()

	nonce := GetNonce(m.superAcc.privateKey)
	fmt.Println("Begin TransferToken0ToAccount", "transferOkT", needTransferToWorker, "nonce", nonce)
	txs := make([]*types.Transaction, 0)
	accCnt := 0
	for _, group := range m.groupList {
		for _, acc := range group.worker {
			if accCnt >= m.paraNum {
				break
			}
			if needTransferToWorker {
				tx := transferOkt(m.superAcc.privateKey, acc.ethAddress, nonce, ether)
				nonce++
				txs = append(txs, tx)
			}

			payload, err := erc20Builder.Build("transfer", acc.ethAddress, new(big.Int).SetInt64(10000000000))
			panicerr(err)
			tx := SendTxWithNonce(m.superAcc.privateKey, group.contract.Token0, payload, nonce)
			nonce++
			txs = append(txs, tx)

			tx = SendTxWithNonce(m.superAcc.privateKey, group.contract.Token2, payload, nonce)
			nonce++
			txs = append(txs, tx)
			accCnt++

		}
	}
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

	var payload []byte
	var err error

	// approve token0
	payload, err = erc20Builder.Build("approve", c.Router, new(big.Int).SetInt64(1000))
	panicerr(err)
	if err = SendTx(a.privateKey, token0, payload); err != nil {
		return err
	}

	// swap token0->token1
	payload, err = routerBuilder.Build("swapExactTokensForTokens",
		new(big.Int).SetInt64(500), new(big.Int),
		[]common.Address{token0, token1},
		a.ethAddress, big.NewInt(1956981781),
	)
	panicerr(err)
	if err = SendTx(a.privateKey, c.Router, payload); err != nil {
		return err
	}

	// approve token0 (for addLiquidity)
	payload, err = erc20Builder.Build("approve", c.Router, new(big.Int).SetInt64(30))
	panicerr(err)
	if err = SendTx(a.privateKey, token0, payload); err != nil {
		return err
	}

	// approve token1 (for addLiquidity)
	payload, err = erc20Builder.Build("approve", c.Router, new(big.Int).SetInt64(30))
	panicerr(err)
	if err = SendTx(a.privateKey, token1, payload); err != nil {
		return err
	}

	// addLiquidity
	payload, err = routerBuilder.Build("addLiquidity", token0, token1, new(big.Int).SetInt64(30), new(big.Int).SetInt64(30), new(big.Int), new(big.Int), a.ethAddress, big.NewInt(1956981781))
	panicerr(err)
	if err = SendTx(a.privateKey, c.Router, payload); err != nil {
		return err
	}

	// approve lp for stakingRewards
	payload, err = erc20Builder.Build("approve", stakeRewards, new(big.Int).SetInt64(10))
	panicerr(err)
	if err = SendTx(a.privateKey, lp, payload); err != nil {
		return err
	}

	// stake lp for stakingRewards
	payload, err = StakingRewardsBuilder.Build("stake", big.NewInt(10))
	panicerr(err)
	if err = SendTx(a.privateKey, stakeRewards, payload); err != nil {
		return err
	}

	// getReward for stakingRewards
	payload, err = StakingRewardsBuilder.Build("getReward")
	panicerr(err)
	if err = SendTx(a.privateKey, stakeRewards, payload); err != nil {
		return err
	}

	payload, err = StakingRewardsBuilder.Build("withdraw", big.NewInt(3))
	panicerr(err)
	if err = SendTx(a.privateKey, stakeRewards, payload); err != nil {
		return err
	}

	payload, err = StakingRewardsBuilder.Build("exit")
	panicerr(err)
	if err = SendTx(a.privateKey, stakeRewards, payload); err != nil {
		return err
	}
	return nil
}

func (m *wmtManager) run(c SwapContract, a *acc) {
	rand.Seed(time.Now().UnixNano())
	sleepTime := rand.Intn(20)
	fmt.Println("run---", c.Token0.String(), a.ethAddress.String(), sleepTime)
	for true {
		time.Sleep(time.Duration(sleepTime) * time.Second)
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
