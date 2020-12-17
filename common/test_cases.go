package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/mitchellh/mapstructure"
)

const (
	// run tx mode
	Serial   = "serial"
	Parallel = "parallel"

	// fee mode
	AUTO = "auto"

	// roun max num
	IntMax = int(^uint(0) >> 1)
)

type TestCases []TestCase

func (cases TestCases) String() string {
	var s bytes.Buffer
	for i, c := range cases {
		s.WriteString(fmt.Sprintf("⚙️ Test Case %d Config ⚙️\n", i))
		s.WriteString(c.String() + "\n")
	}
	return s.String()
}

type TestCase struct {
	MnemonicPath string `json:"mnemonic_path"`
	RunTxMode    string `json:"run_tx_mode"`
	BaseParam
	Transactions []Transaction `json:"transactions"`
}

func (c TestCase) String() string {
	var s bytes.Buffer
	if c.MnemonicPath == "" {
		s.WriteString("MnemonicPath:     nil " + "(test accounts will be generated automatically)" + "\n")
	} else {
		s.WriteString("MnemonicPath:     " + c.MnemonicPath + "\n")
	}
	s.WriteString("RunTxMode:        " + c.RunTxMode + "\n")
	s.WriteString("Fee:              " + c.Fee + "\n")
	s.WriteString("ConcurrentNum:    " + strconv.Itoa(c.ConcurrentNum) + "\n")
	s.WriteString("SleepTime:        " + strconv.Itoa(c.SleepTime) + "s" + "\n")
	if c.RoundNum == IntMax {
		s.WriteString("Round:            " + "∞" + "\n")
	} else {
		s.WriteString("Round:            " + strconv.Itoa(c.RoundNum) + "\n")
	}

	for i, tx := range c.Transactions {
		s.WriteString(fmt.Sprintf("    🐯🐯 Transaction %d Config 🐯🐯\n", i))
		s.WriteString("    Type:             " + tx.Type + "\n")
		var arg BaseParam
		err := mapstructure.Decode(tx.Args, &arg)
		if err != nil {
			log.Fatalf("failed to decode args config. error: %s\n", err.Error())
		}
		arg.SetBaseParam(c.BaseParam)
		s.WriteString("    Fee:              " + arg.Fee + "\n")
		s.WriteString("    ConcurrentNum:    " + strconv.Itoa(arg.ConcurrentNum) + "\n")
		s.WriteString("    SleepTime:        " + strconv.Itoa(arg.SleepTime) + "s" + "\n")
		if arg.RoundNum == IntMax {
			s.WriteString("    Round:            " + "∞" + "\n")
		} else {
			s.WriteString("    Round:            " + strconv.Itoa(arg.RoundNum) + "\n")
		}
	}
	return s.String()
}

// Transactions used for add multi types tx
type Transaction struct {
	Type string                 `json:"type"`
	Args map[string]interface{} `json:"args"`
}

// BaseParm used for decode from the arg map
// every specific param have to extend this structure
type BaseParam struct {
	Fee           string `json:"fee" mapstructure:"fee"`
	ConcurrentNum int    `json:"concurrent_num" mapstructure:"concurrent_num"`
	SleepTime     int    `json:"sleep_time" mapstructure:"sleep_time"`
	RoundNum      int    `json:"round_num" mapstructure:"round_num"`
}

func (p *BaseParam) SetBaseParam(baseParam BaseParam) {
	if p.Fee == "" {
		p.Fee = baseParam.Fee
	}
	if p.SleepTime <= 0 {
		p.SleepTime = baseParam.SleepTime
	}
	if p.RoundNum <= 0 {
		p.RoundNum = baseParam.RoundNum
	}
	if p.ConcurrentNum <= 0 {
		p.ConcurrentNum = baseParam.ConcurrentNum
	}
}

func ReadTestCases(txConfigPath string) (TestCases, error) {
	data, err := ioutil.ReadFile(txConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load tx json : %s", err.Error())
	}

	var testcases TestCases
	err = json.Unmarshal(data, &testcases)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal error : %s", err.Error())
	}

	for i := range testcases {
		validate(&testcases[i])
	}
	return testcases, nil
}

func validate(c *TestCase) {
	if c.RunTxMode == "" {
		c.RunTxMode = Parallel
	}
	if c.Fee == "" {
		c.Fee = AUTO
	} else {
		_, err := types.NewDecFromStr(c.Fee)
		if err != nil {
			log.Fatalln("failed when validate test case. error:", err.Error())
		}
	}
	if c.ConcurrentNum <= 0 {
		c.ConcurrentNum = 1
	}
	if c.SleepTime <= 0 {
		c.SleepTime = 5
	}
	if c.RoundNum <= 0 {
		c.RoundNum = IntMax
	}
}
