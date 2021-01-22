package keeper

// to vote for all target vals
func (k *Keeper) RaisePercentageToPlunder() error {
	// choose a pure worker
	worker, err := k.getPureWorker()
	if err != nil {
		return err
	}

	// vote for all target validators
	return k.addSharesToAllVals(worker)
}
