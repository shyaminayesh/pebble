package seed

import (
	"fmt"
	"strings"
	"io/ioutil"
	"pebble/config"
	"pebble/utils/log"
	"github.com/spf13/cobra"
)

var Seed_Command = &cobra.Command{
	Use: "seed",
	Short: "Seed database",
	Run: seed,
}

func seed(cmd *cobra.Command, args []string) {

	/*
		We need to query configuration information from
		the main configuration file to continue the 
		process
	*/
	var configs = config.Config()
	seed_configs := configs.Sub("seed")


	/*
		We have to read all the seed files and then
		act according to the file instructions.
	*/
	var yamls []string
	files, err := ioutil.ReadDir("./" + seed_configs.Get("dir").(string))
	if err != nil {
		logger.Println("red", "[FATAL]: ", "Failed to load seed files.")
	}
	for _, file := range files {
		if strings.Contains(file.Name(), ".yml") {
			yamls = append(yamls, strings.Replace(file.Name(), ".yml", "", -1))
		}
	}


	// print
	fmt.Println( yamls )
}