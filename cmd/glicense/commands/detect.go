package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	paramDetectPath string
	paramDetectURL  string
)

func init() {
	detectCmd.Flags().StringVarP(&paramDetectPath, "path", "p", "", "Path of license file")
	detectCmd.Flags().StringVarP(&paramDetectURL, "url", "u", "", "Url path to license file")
	rootCmd.AddCommand(detectCmd)
}

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect license info based on license content",
	Long: `
Detect license info based on license content.
The inputted license content can be not entirely, but we will find matchest license with your input.

Support 3 ways to use:
	glicense detect "MIT License Copyright (c) Permission is hereby granted..."
	glicense detect -p /path/to/source/
	glicense detect -u https://github.com/abc/

We also support pipe input :
	echo "MIT License Copyright (c) Permission is hereby granted..." | glicense detect

Note:
 - Direct license content will have a higher priority to get.
 - If the detect command is set both arguments of a local path (-p) and URL (-u), local path will be used.

`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Detect function - under construction")
	},
}
