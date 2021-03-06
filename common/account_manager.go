package common

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/okexchain-go-sdk/utils"
)

type AccountManager struct {
	i     int
	infos []keys.Info
	sum   int
	lock  *sync.Mutex
}

func newAccountManager(infos []keys.Info) *AccountManager {
	return &AccountManager{
		0,
		infos,
		len(infos),
		new(sync.Mutex),
	}
}

func (m *AccountManager) GetInfo() keys.Info {
	m.lock.Lock()
	defer m.lock.Unlock()
	k := m.i
	m.i = (m.i + 1) % m.sum
	return m.infos[k]
}

func (m *AccountManager) GetInfos() []keys.Info {
	return m.infos
}

func GetAccountManagerFromFile(path string, limit ...int) *AccountManager {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalln(err.Error())
		return nil
	}
	defer f.Close()
	fmt.Printf("loading mnemonic %s, please wait\n", path)

	num := 9999999999 // BIG NUMBER
	if len(limit) != 0 {
		num = limit[0]
	}

	var infos []keys.Info
	rd := bufio.NewReader(f)
	for index := 0; index < num; index++ {
		mnemonic, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}

		info, _, err := utils.CreateAccountWithMnemo(strings.TrimSpace(mnemonic), fmt.Sprintf("acc%d", index), PassWord)
		if err != nil {
			log.Fatalln(err.Error())
		}
		infos = append(infos, info)
		//fmt.Println(info.GetAddress().String(), index)
	}
	return newAccountManager(infos)
}

func GetAccountAddressFromFile(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalln(err.Error())
		return nil
	}
	defer f.Close()

	fmt.Printf("loading address %s, please wait\n", path)
	var addrs []string
	rd := bufio.NewReader(f)
	for {
		addr, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}

		addrs = append(addrs, strings.TrimSpace(addr))
	}
	return addrs
}

func GetPrivKeyFromMnemoFile(path string, limit ...int) (privKeys []string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalln(err.Error())
		return nil
	}
	defer f.Close()

	num := 9999999999 // BIG NUMBER
	if len(limit) != 0 {
		num = limit[0]
	}


	fmt.Printf("loading mnemonics from path: %s, please wait\n", path)
	rd := bufio.NewReader(f)
	for index := 0; index < num; index++ {
		mnemo, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}

		privKey, err := utils.GeneratePrivateKeyFromMnemo(strings.TrimSpace(mnemo))
		if err != nil {
			panic(err)
		}
		fmt.Println(privKey)

		privKeys = append(privKeys, privKey)
	}

	return
}

func GetPrivKeyFromPrivKeyFile(path string, limit ...int) (privKeys []string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalln(err.Error())
		return nil
	}
	defer f.Close()

	num := 9999999999 // BIG NUMBER
	if len(limit) != 0 {
		num = limit[0]
	}

	fmt.Printf("loading privkey from path: %s, please wait\n", path)
	rd := bufio.NewReader(f)
	for index := 0; index < num; index++ {
		privKey, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}

		privKeys = append(privKeys, strings.TrimSpace(privKey))
	}
	return
}

func GetAccountAddressFromMnemoFile(path string) (addrs []string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalln(err.Error())
		return nil
	}
	defer f.Close()

	fmt.Printf("loading mnemonics from path: %s, please wait\n", path)
	rd := bufio.NewReader(f)
	for {
		mnemo, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}

		info, _, err := utils.CreateAccountWithMnemo(strings.TrimSpace(mnemo), "acc", PassWord)
		if err != nil {
			panic(err)
		}

		addr := info.GetAddress().String()
		fmt.Println(addr)
		addrs = append(addrs, addr)
	}

	return
}
