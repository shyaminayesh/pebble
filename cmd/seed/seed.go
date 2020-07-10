package seed

import (
	"os"
	"fmt"
	"strings"
	"io/ioutil"
	"database/sql"
	"pebble/config"
	"pebble/utils/log"
	"github.com/spf13/cobra"
	_ "github.com/go-sql-driver/mysql"
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
	connection_configs := configs.Sub("connection")
	seed_configs := configs.Sub("seed")


	dialect := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", connection_configs.Get("user"), connection_configs.Get("password"), connection_configs.Get("host"), connection_configs.Get("port"), connection_configs.Get("name"))
	db, err := sql.Open(connection_configs.Get("driver").(string), dialect)
	if err != nil {
		logger.Println("red", "[FATAL]: ", "Database connection failed.")
		os.Exit(1)
	}
	defer db.Close()


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


	/*
		Iterate through each yaml file to apply
		the changes one by one to the database.
	*/
	for _, yaml := range yamls {

		/*
			Check if the database table is already exists with the
			provided yaml file name before inserting any records.
			This will avoid popping errors by the database.
		*/
		var table_count int
		query := fmt.Sprintf("SELECT CAST(COUNT(TABLE_NAME) AS UNSIGNED) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA='%s' AND TABLE_NAME='%s'", connection_configs.Get("name"), yaml)
		db.QueryRow(query).Scan(&table_count)
		if table_count == 0 {
			logger.Println("red", "[FATAL]: ", "TABLE NOT FOUND ( " + yaml + " )")
			os.Exit(1)
		}

		fmt.Println( yaml )
	}

	// print
	fmt.Println( yamls )
}