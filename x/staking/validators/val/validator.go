package val

import (
	"sync"

	"github.com/okex/exchain-go-sdk"
)

type Validator interface {
	Create(wg *sync.WaitGroup)
	Edit(wg *sync.WaitGroup)
	Destroy(wg *sync.WaitGroup)
}

// validator will invoke gosdk to submit create-validator or destory-validator tx
type validator struct {
	state  int // bonded, unbonding, unboned
	jailed bool
	c      *gosdk.Client
}

func (v *validator) Create() {

}

func (v *validator) Destroy() {

}
