/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// seekCmd represents the seek command
var seekCmd = &cobra.Command{
	Use:   "seek",
	Short: "A brief description of your command",
	Long: ``,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("seek called")
	},
}

func init() {
	rootCmd.AddCommand(seekCmd)
}
