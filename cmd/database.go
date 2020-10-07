package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Database_Command = &cobra.Command{
	Use: "database",
	Short: "Manage database",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: pebble database [OPTION]")
		fmt.Println()
		fmt.Println("   database")
		fmt.Println("       This command will migrate all the")
		fmt.Println("       tables defined in the schema directory")
		fmt.Println()
	},
}
