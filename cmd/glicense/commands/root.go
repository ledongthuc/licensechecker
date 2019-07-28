package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "glicense",
	Short: "G-License is toolkit for detect license and add license content to source files",
	Long: `
A tool kit for working with license.
 - Listing license info and content.
 - Finding license info and content by name.
 - Detect license info by their content.
 - Add license content into your source code files.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
