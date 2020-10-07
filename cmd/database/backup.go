package database

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Database_Backup_Command = &cobra.Command{
	Use: "backup",
	Short: "Backup database",
	Run: run_backup,
}


func run_backup(cmd *cobra.Command, args []string) {
	fmt.Println("run_backup")
}