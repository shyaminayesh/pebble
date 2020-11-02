package schema

import (
	"fmt"
	"log"
	"bytes"
	"strings"
	"io/ioutil"
	"database/sql"
	"pebble/config"
	"pebble/utils/log"
	"github.com/spf13/cobra"
	parser "pebble/utils/parser/schema"
	schemalex "github.com/schemalex/schemalex/diff"
)

var Schema_Migrate_Command = &cobra.Command{
	Use: "migrate",
	Short: "Migrate schema",
	Run: schema_migrate,
}

func schema_migrate(cmd *cobra.Command, args []string) {

	/**
	* We have to load main pebble configuration to get
	* essential configuration information to continue
	* the process.
	*/
	Config := config.Get()

	/**
	* Establish connection to the database using the
	* pebble configuration.
	*/
	dialect := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", Config.Connection.User, Config.Connection.Password, Config.Connection.Host, Config.Connection.Port, Config.Connection.Database)
	db, err := sql.Open(Config.Connection.Dialect, dialect)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	/*
		We need to iterate through every file in the
		migration directory to find every migration answer
		file to operate.
	*/
	var schemas []string
	files, err := ioutil.ReadDir("./" + Config.Schema.Directory)
	if err != nil { log.Fatal(err) }

	for _, file := range files {
		if strings.Contains(file.Name(), ".yml") {
			schemas = append(schemas, strings.Replace(file.Name(), ".yml", "", -1))
		}
	}

	/*
		We need to drop tables that are not in our schema
		set first to clear out the unwanted tables from
		the database.
	*/
	query := fmt.Sprintf("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA='%s' AND TABLE_NAME NOT IN ('%s')", Config.Connection.Database, strings.Join(schemas, "','"))
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	/*
		We have the names of the tables that has no schema
		file present in the migration store. We can safely
		drop those tables from the database.
	*/
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			log.Fatal(err)
		}
		logger.Println("green", "[TABLE]: ", "DROPPING ( " + table + " )")
		db.Exec("DROP TABLE " + table)
	}



	// HERE ####
	for _, schema := range schemas {
		fmt.Println("[TABLE]: " + schema)

		/**
		* Check if the table is exists or not and then create
		* the table if it's not exists first.
		*/
		var count int
		query := fmt.Sprintf("SELECT CAST(COUNT(TABLE_NAME) AS UNSIGNED) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA='%s' AND TABLE_NAME='%s'", Config.Connection.Database, schema)
		db.QueryRow(query).Scan(&count)


		/**
		* Here we handle the table not exists state by using the
		* table count from the last sql query and then we create new
		* table in the database according to the migration file.
		*/
		if count == 0 {
			parser := parser.Schema {}
			parser.File("./" + Config.Schema.Directory + "/" + schema + ".yml")
			db.Exec(parser.Statement())
		}

		/**
		* Depending on the recent sql query table count details we can
		* decide to modify exsisting table schema according to the
		* changes in the migration file.
		*/
		if count >= 1 {

			parser := parser.Schema {}
			parser.File("./" + Config.Schema.Directory + "/" + schema + ".yml")

			result := []string{"table", "ddl"}
			db.QueryRow(fmt.Sprintf("SHOW CREATE TABLE %s", schema)).Scan(&result[0], &result[1])

			statement := &bytes.Buffer{}
			err := schemalex.Strings(statement, result[1], parser.Statement())
			if err != nil { log.Fatal(err) }

			for _, stmnt := range strings.Split(statement.String(), ";") {
				db.Exec(stmnt)
			}

		}

	}

}
