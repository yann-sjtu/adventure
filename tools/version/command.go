package version

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Cmd prints out the application's version information passed via build flags.
var Cmd = &cobra.Command{
	Use:   "version",
	Short: "Print the advanture version",
	RunE: func(_ *cobra.Command, _ []string) error {
		verInfo := NewInfo()
		bz, err := yaml.Marshal(&verInfo)
		if err != nil {
			return err
		}

		_, err = fmt.Printf(string(bz))
		return err
	},
}
