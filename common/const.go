package common

import "github.com/okex/okexchain/x/common"

const (
	PassWord          = "12345678"
	NativeToken       = common.NativeToken
	DefaultStableCoin = "usdk"
	RichMnemonic      = "actual assume crew creek furnace water electric fitness stumble usage embark ancient"
)

const (

	// distribution
	WithdrawRewards = "withdraw-rewards"
	SetWithdrawAddr = "set-withdraw-addr"

	//token
	Issue                  = "issue"
	Burn                   = "burn"
	Mint                   = "mint"
	MultiSend              = "multi-send"
	TokenTransferOwnership = "token-transfer-ownership"
	Edit                   = "edit"

	//dex
	List                 = "list"
	Deposit              = "deposit"
	Withdraw             = "withdraw"
	DexTransferOwnership = "dex-transfer-ownership"
	RegisterOperator     = "register-operator"
	EditOperator         = "edit-operator"

	//order
	Order = "order"

	//staking
	DelegateVoteUnbond = "delegate_vote_unbond"
	Proxy              = "proxy"

	//ammswap
	AddLiquidity    = "add-liquidity"
	RemoveLiquidity = "remove-liquidity"
	CreateExchange  = "create-exchange"
	SwapExchange    = "swap-exchange"
)
