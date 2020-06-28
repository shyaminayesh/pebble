package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Version_Command = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v1.0.0-beta")
	},
}