package locker

import (
	"fmt"
	"sync"

	"github.com/okex/adventure/x/strategy/farm/client"
	"github.com/okex/adventure/x/strategy/farm/emitter"
	poolertypes "github.com/okex/adventure/x/strategy/farm/pooler/types"
	"github.com/okex/adventure/x/strategy/farm/utils"
	"github.com/spf13/cobra"
)

const (
	flagPoolerFilePath     = "pooler-path"
	flagLockerAddrFilePath = "locker-path"

	// params
	issuedTokenAmount = 100
	lptTokenAmount    = 1
	addrNumOneTime    = 50
)

func allocateTokensToAllLockersFromAllPoolersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "allocate-tokens-to-lockers-from-all-poolers",
		Short: "allocate all tokens to all lockers from all poolers",
		Args:  cobra.NoArgs,
		RunE:  runAllocateTokensToAllLockersFromAllPoolersCmd,
	}

	flags := cmd.Flags()
	flags.StringP(flagPoolerFilePath, "p", "", "the file path of pooler mnemonics")
	flags.StringP(flagLockerAddrFilePath, "l", "", "the file path of locker addresses")

	return cmd
}

func runAllocateTokensToAllLockersFromAllPoolersCmd(cmd *cobra.Command, args []string) error {
	// load pooler manager
	poolerPath, err := cmd.Flags().GetString(flagPoolerFilePath)
	if err != nil {
		return err
	}

	emt := emitter.NewEmitter(poolertypes.GetPoolerManager(poolerPath), nil)

	// get locker address
	lockerAddrsPath, err := cmd.Flags().GetString(flagLockerAddrFilePath)
	if err != nil {
		return err
	}

	lockerAddrsStrs, err := utils.GetTestAddrsFromFile(lockerAddrsPath)
	if err != nil {
		panic(err)
	}

	lockerAddrs, err := utils.ParseAccAddrsFromAddrsStr(lockerAddrsStrs)
	if err != nil {
		panic(err)
	}

	for _, pooler := range emt.PoolerManager {
		cli := client.CliManager.GetClient()
		if pooler.UpdateIssueTokenInfo(cli, true) != nil {
			panic(err)
		}
	}

	var wg sync.WaitGroup
	for _, pooler := range emt.PoolerManager {
		wg.Add(1)
		go func(pPooler *poolertypes.Pooler) {
			if err := utils.MultiSendByGroup(&wg, pPooler.GetKey(), fmt.Sprintf("%d%s", issuedTokenAmount, pPooler.GetIssuedTokenSymbol()), lockerAddrs, addrNumOneTime); err != nil {
				fmt.Println(err)
			}
		}(pooler)
	}
	wg.Wait()

	for _, pooler := range emt.PoolerManager {
		wg.Add(1)
		go func(pPooler *poolertypes.Pooler) {
			if err := utils.MultiSendByGroup(&wg, pPooler.GetKey(), fmt.Sprintf("%d%s", lptTokenAmount, pPooler.GetLPTSymbol()), lockerAddrs, addrNumOneTime); err != nil {
				fmt.Println(err)
			}
		}(pooler)
	}
	wg.Wait()

	return nil
}
