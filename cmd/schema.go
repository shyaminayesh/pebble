package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Schema_Command = &cobra.Command{
	Use: "schema",
	Short: "Manage DB schema",
	Run: run,
}

var config = Config()

func run(cmd *cobra.Command, args []string) {

	connection := config.Sub("connection")
	fmt.Println( connection.Get("port") )

	for _, arg := range args {
		fmt.Println( arg )
	}

}