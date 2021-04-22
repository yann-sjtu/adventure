package validators

import (
	"log"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	stakingTypes "github.com/okex/exchain/x/staking/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmTypes "github.com/tendermint/tendermint/types"
)

var (
	valNumber = 0
)

func createValidatorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validators",
		Short: "create fake validators on testnet",
		Args:  cobra.NoArgs,
		RunE:  runCreateValidators,
	}
	flags := cmd.Flags()
	flags.IntVarP(&valNumber, "number", "n", 1, "the number of the fake validators")
	flags.StringP(flagAccountsFilePath, "p", "", "the file path of validators mnemonic")
	return cmd
}

//okchaincli tx staking create-validator --pubkey=$(okchaind tendermint show-validator)
// --moniker="my nickname" --identity="logo|||http://mywebsite/pic/logo.jpg" --website="http://mywebsite" --details="my slogan"
// --from jack

type fakeValidator struct {
	Address tmTypes.Address `json:"address"`
	PubKey  crypto.PubKey   `json:"pub_key"`
	PrivKey crypto.PrivKey  `json:"priv_key"`
}

func runCreateValidators(cmd *cobra.Command, args []string) error {
	path, err := cmd.Flags().GetString(flagAccountsFilePath)
	if err != nil {
		return err
	}
	accManager := common.GetAccountManagerFromFile(path, valNumber)

	clientManager := common.NewClientManager(common.Cfg.Hosts, common.AUTO)
	for i := 0; i < valNumber; i++ {
		privKey := ed25519.GenPrivKey()
		fakeValidator := &fakeValidator{
			Address: privKey.PubKey().Address(),
			PubKey:  privKey.PubKey(),
			PrivKey: privKey,
		}
		pubkey, err := stakingTypes.Bech32ifyConsPub(fakeValidator.PubKey)
		if err != nil {
			return err
		}

		acc := accManager.GetInfo()
		cli := clientManager.GetClient()
		accInfo, err := cli.Auth().QueryAccount(acc.GetAddress().String())
		if err != nil {
			log.Printf("[%d] failed. val %s query before create: %s\n", i, sdk.ValAddress(acc.GetAddress()).String(), err.Error())
			return err
		}

		if _, err := cli.Staking().CreateValidator(acc, common.PassWord, pubkey,
			"test"+strconv.Itoa(i), "", "", "",
			"", accInfo.GetAccountNumber(), accInfo.GetSequence()); err != nil {
			log.Printf("[%d] failed. val %s try to be created: %s\n", i, sdk.ValAddress(acc.GetAddress()).String(), err.Error())
			return err
		}
		log.Printf("[%d] %s create successfully\n", i, pubkey)
	}

	return nil
}
