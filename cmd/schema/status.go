package schema

import (
	"fmt"
	"pebble/config"
	"github.com/spf13/cobra"
)

var Schema_Status_Command = &cobra.Command{
	Use: "status",
	Short: "View schema status",
	Run: run,
}

var conf = config.Config()

func run(cmd *cobra.Command, args []string) {
	fmt.Println("schema::status")
}