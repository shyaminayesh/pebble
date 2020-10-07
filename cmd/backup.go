package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Backup_Command = &cobra.Command{
	Use: "backup",
	Short: "Manage database backups",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: pebble backup [OPTION]")
		fmt.Println()
		fmt.Println("   This backup sub command will help you to")
		fmt.Println("   manage databse backups using the pebble ")
		fmt.Println("   command line                            ")
		fmt.Println()
		fmt.Println("   create                                  ")
		fmt.Println("       create command will create a backup ")
		fmt.Println("       of your database according to the   ")
		fmt.Println("       params you defined in your main     ")
		fmt.Println("       pebble configuration file           ")
		fmt.Println()
	},
}
