package seed

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Seed_Command = &cobra.Command{
	Use: "seed",
	Short: "Seed database",
	Run: seed,
}

func seed(cmd *cobra.Command, args []string) {



	fmt.Println("seed")
}