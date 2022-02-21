package scenario

import "github.com/spf13/cobra"

func ScenarioCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "Scenario",
		Short: "Scenario test including getBalance, tx and getBalance",
		Run: scenario,
	}
	cmd.Flags().BoolVarP(&fixed, "fixed", "f", false, "if true, transfer to one address; otherwise, transfer to a random address")
	return cmd
}

