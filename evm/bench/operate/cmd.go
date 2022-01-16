package operate

import "github.com/spf13/cobra"

func OperateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operate",
		Short: "before execute, it requires that user deploys contracts from github.com/okex/evm-performance",
		Run:   operate,
	}

	cmd.Flags().BoolVar(&direct, "direct", false, "if true, contract address should be Router; otherwise, should be Test-Contract")
	cmd.Flags().StringVar(&contract, "contract", "", "contract which implements IOperate interface, or Router Address")
	cmd.Flags().Int64Var(&id, "id", 0, "Test-Contract id in Router")

	cmd.Flags().Int64SliceVar(&opts, "opts", []int64{}, "")
	cmd.Flags().Int64Var(&times, "times", 0, "")
	return cmd
}
