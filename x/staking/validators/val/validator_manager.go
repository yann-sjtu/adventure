package val

import "time"

// manages all validators
type ValManager struct {
	vset []*validator // 40 validators
	run  bool
}

// use gosdk to get all validators state
func (m *ValManager) getState() {
	// do what 'okchaincli query staking validators' does
}

// run in a dedicated thread
func (m *ValManager) Housekeeping() {

	for {
		m.getState()

		// get all validators state and submit create or destroy tx accordingly

		// ...

		time.Sleep(time.Minute * 1)

		if !m.run {
			break
		}
	}
}

func (m *ValManager) Stop() {
	m.run = false
}
