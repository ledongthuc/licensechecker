package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	paramAddPath   string
	paramAddConfig string
)

func init() {
	addCmd.Flags().StringVarP(&paramAddPath, "path", "p", "", "Path of license content file")
	addCmd.Flags().StringVarP(&paramAddConfig, "config", "c", "", "Config file for mapping/excluded files")
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add license content into source code files.",
	Long: `
Add license content into source code files. It will look like:
	/* License content ... */
	...Source code
	/* License content ... */

Usage:
	glicense add "MIT License Copyright (c) Permission is hereby granted..." /path/to/source/code/to/add/
	glicense add -p file/license.txt /path/to/source/code/to/add/
	glicense add -c config/file.conf -p file/license.txt  /path/to/source/code/to/add/

Config file:
 - You want support other source file format to add license into?
 - You want to use your custom pattern when adding license content into extensions?
 - You want to exclude some file/path with pattern from adding license conten?
 Config file can control it with simple format: {
	excluded_pattern: "*.sql",
	mapping_comment: {
		".go": "/*{content}*/",
		".xml": "<!--{content}-->",
		".custom": "/*{content}*/",
	},
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Add function - under construction")
	},
}
