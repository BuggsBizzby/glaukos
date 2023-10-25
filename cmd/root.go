/*
Copyright © 2023 RIPSFORFUN

*/
package cmd

import (
	"os"
	"github.com/spf13/cobra"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "glaukos",
	Short: "",
	Long: `Glaukos: Lord of the Phishermen
    A multi-environment creation tool, utilizing KasmVNC for remote viewing and mitmproxy for data collection.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
    rootCmd.CompletionOptions.HiddenDefaultCmd = true
}


