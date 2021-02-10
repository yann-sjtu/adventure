package UniswapV2Staker

import "github.com/okex/okexchain-go-sdk/utils"

var (
	StakingRewardsBuilder utils.PayloadBuilder
	ERC20Builder          utils.PayloadBuilder
)

func Init() {
	var err error

	// 1. init builders
	StakingRewardsBuilder, err = utils.NewPayloadBuilder(StakingRewardsBin, StakingRewardsABI)
	if err != nil {
		panic(err)
	}

	ERC20Builder, err = utils.NewPayloadBuilder(ERC20Bin, ERC20ABI)
	if err != nil {
		panic(err)
	}
}

func BuildExitPayload() []byte {
	payload, err := StakingRewardsBuilder.Build("exit")
	if err != nil {
		panic(err)
	}
	return payload
}

func BuildGetRewardPayload() []byte {
	payload, err := StakingRewardsBuilder.Build("getReward")
	if err != nil {
		panic(err)
	}
	return payload
}

func BuildStakePayload(num int) []byte {
	payload, err := StakingRewardsBuilder.Build("stake", utils.Uint256(num))
	if err != nil {
		panic(err)
	}
	return payload
}

func BuildWithdrawPayload(num int) []byte {
	payload, err := StakingRewardsBuilder.Build("withdraw", utils.Uint256(num))
	if err != nil {
		panic(err)
	}
	return payload
}

func BuildApprovePayload(addr string, amount int) []byte {
	payload, err := ERC20Builder.Build("approve", utils.EthAddress(addr), utils.Uint256(amount))
	if err != nil {
		panic(err)
	}
	return payload
}
