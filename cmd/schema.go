package cmd

import (
	"github.com/spf13/cobra"
)

var Schema_Command = &cobra.Command{
	Use: "schema",
	Short: "Manage database schema",
	Run: run,
}

func run(cmd *cobra.Command, args []string) {
	// for _, arg := range args {
	// 	fmt.Println( arg )
	// }
}