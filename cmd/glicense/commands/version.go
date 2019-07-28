package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VERSION = "0.1.0"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of G-License application",
	Long:  "Print the version number of G-License application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("G-License", VERSION)
	},
}
