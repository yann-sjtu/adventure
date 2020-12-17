package gov

import (
	"fmt"
	"io/ioutil"

	"github.com/okex/adventure/common"
)

var (
	textProposalFilePath               = "./text_proposal.json"
	paramChangeProposalFilePath        = "./param_change_proposal.json"
	delistProposalFilePath             = "./delist_proposal.json"
	communityPoolSpendProposalFilePath = "./community_pool_spend_proposal.json"

	textProposalJSON               = fmt.Sprintf(`{"title":"Text Proposal","description":"text","proposal_type":"Text","deposit":"100%s"}`, common.NativeToken)
	paramChangeProposalJSON        = fmt.Sprintf(`{"title":"Param Change Proposal","description":"param change proposal description","changes":[{"subspace":"staking","key":"MaxValidators","value":105}],"deposit":[{"denom":"%s","amount":"100"}],"height":"1024"}`, common.NativeToken)
	delistProposalJSON             = fmt.Sprintf(`{"title":"Delist Proposal","description":"delist proposal description","base_asset":"btc-000","quote_asset":"%s","deposit":[{"denom":"%s","amount":"100"}]}`, common.NativeToken, common.NativeToken)
	communityPoolSpendProposalJSON = fmt.Sprintf(`{"title":"Community Pool Spend Proposal","description":"community pool spend description","recipient":"okchain1hw4r48aww06ldrfeuq2v438ujnl6alszzzqpph","amount":[{"denom":"%s","amount":"10.24"}],"deposit":[{"denom":"%s","amount":"100"}]}`, common.NativeToken, common.NativeToken)
)

func writeProposal() {
	_ = ioutil.WriteFile(communityPoolSpendProposalFilePath, []byte(communityPoolSpendProposalJSON), 0644)
}
