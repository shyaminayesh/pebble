package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Schema_Command = &cobra.Command{
	Use: "schema",
	Short: "Manage database schema",
	Run: run,
}

func run(cmd *cobra.Command, args []string) {

	fmt.Println("Usage: pebble schema [OPTION]")
	fmt.Println()
	fmt.Println("   migrate")
	fmt.Println("       This command will migrate all the")
	fmt.Println("       tables defined in the schema directory")
	fmt.Println()

}