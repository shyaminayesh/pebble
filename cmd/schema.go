package cmd

import (
	"fmt"
	"pebble/config"
	"github.com/spf13/cobra"
)

var Schema_Command = &cobra.Command{
	Use: "schema",
	Short: "Manage database schema",
	Run: run,
}

var conf = config.Config()

func run(cmd *cobra.Command, args []string) {

	connection := conf.Sub("connection")
	fmt.Println( connection.Get("port") )

	for _, arg := range args {
		fmt.Println( arg )
	}

}