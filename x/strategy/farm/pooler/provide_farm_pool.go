package pooler

//func provideFarmPoolCmd() *cobra.Command {
//	cmd := &cobra.Command{
//		Use:   "provide-farm-pool",
//		Short: "provide farm pool that created by pooler himself",
//		Args:  cobra.NoArgs,
//		RunE:  runProvideFarmPool,
//	}
//
//	flags := cmd.Flags()
//	flags.StringP(flagIssuerFilePath, "p", "", "the file path of pooler mnemonics")
//
//	return cmd
//}

//func runProvideFarmPool(cmd *cobra.Command, args []string) error {
//	path, err := cmd.Flags().GetString(flagIssuerFilePath)
//	if err != nil {
//		return err
//	}
//	poolerManager := types.GetPoolerManagerFromFiles(path, "pooler")
//
//	pools, err := client.CliManager.GetClient().Farm().QueryPools()
//	if err != nil {
//		return err
//	}
//
//	currentHeight, err := utils.QueryLatestBlockHeight()
//	if err != nil {
//		return err
//	}
//
//	groupNum := 100
//	times := len(pools)/groupNum + 1
//	for i := 0; i < times; i++ {
//		var index2 int
//		index1 := i * groupNum
//		if i != times-1 {
//			index2 = (i + 1) * groupNum
//		} else {
//			index2 = len(pools)
//		}
//		runProvideFarmPoolsParallelly(pools[index1:index2], currentHeight, poolerManager)
//
//		fmt.Printf("Group %d: provide pool %d ~ %d ...\n", i, index1, index2)
//		time.Sleep(5 * time.Second)
//	}
//
//	return nil
//}
//
//func runProvideFarmPoolsParallelly(pools []sdk.FarmPool, currentHeight int, poolerManager types.PoolerManager) {
//	// provide farm pools
//	var wg sync.WaitGroup
//	for _, pool := range pools {
//		wg.Add(1)
//		go poolerManager[pool.Owner.String()].ProvideFarmPool(&wg, pool, currentHeight)
//	}
//	wg.Wait()
//}
