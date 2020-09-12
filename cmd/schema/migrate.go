package schema

import (
	"fmt"
	"log"
	"bytes"
	"strings"
	"io/ioutil"
	"database/sql"
	"pebble/utils/log"
	"pebble/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	parser "pebble/utils/parser/schema"
	schemalex "github.com/schemalex/schemalex/diff"
	_ "github.com/go-sql-driver/mysql"
)

var Schema_Migrate_Command = &cobra.Command{
	Use: "migrate",
	Short: "Migrate schema",
	Run: schema_migrate,
}

func schema_migrate(cmd *cobra.Command, args []string) {

	/*
		We need to initialize configuration instance to
		query main configuration file for future usage
	*/
	var conf = config.Config()
	conf_connection := conf.Sub("connection")
	conf_schema := conf.Sub("schema")

	/*
		We need to establish database connection to
		perform queries on the database
	*/
	dialect := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", conf_connection.Get("user"), conf_connection.Get("password"), conf_connection.Get("host"), conf_connection.Get("port"), conf_connection.Get("name"))
	db, err := sql.Open(conf_connection.Get("driver").(string), dialect)
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
	files, err := ioutil.ReadDir("./" + conf_schema.Get("dir").(string))
	if err != nil {
		log.Fatal(err)
	}

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
	query := fmt.Sprintf("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA='%s' AND TABLE_NAME NOT IN ('%s')", conf_connection.Get("name"), strings.Join(schemas, "','"))
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

		type (
			Column struct {
				Name		string		`mapstructure:"name"`
				Type		string		`mapstructure:"type"`
				Nullable	bool		`mapstructure:"nullable"`
				Primary		bool		`mapstructure:"primary"`
				Increment	bool		`mapstructure:"increment"`
				Collation	string		`mapstructure:"collation"`
			}
			Table struct {
				Engine		string		`mapstructure:"engine"`
				Charset		string		`mapstructure:"charset"`
				Collation	string		`mapstructure:"collation"`
			}
			Structure struct {
				Table		Table		`mapstructure:"table"`
				Columns		[]Column	`mapstructure:"columns"`
			}
		)

		v := viper.New()
		v.SetConfigName(schema)
		v.SetConfigType("yml")
		v.AddConfigPath("./" + conf_schema.Get("dir").(string))
		err := v.ReadInConfig()
		if err != nil {
			log.Fatal(err)
		}
		var structure Structure
		v.Unmarshal(&structure)

		/*
			Check if the table is exists or not and then create
			the table if it's not exists first.
		*/
		var count int
		query := fmt.Sprintf("SELECT CAST(COUNT(TABLE_NAME) AS UNSIGNED) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA='%s' AND TABLE_NAME='%s'", conf_connection.Get("name"), schema)
		db.QueryRow(query).Scan(&count)


		/*
			Here we handle the table not exists state by using the
			table count from the last sql query and then we create new
			table in the database according to the migration file.
		*/
		if count == 0 {
			parser := parser.Schema {}
			parser.File("./" + conf_schema.Get("dir").(string) + "/" + schema + ".yml")
			db.Exec(parser.Statement())
		}

		/*
			Depending on the recent sql query table count details we can
			decide to modify exsisting table schema according to the
			changes in the migration file.
		*/
		if count >= 1 {

			parser := parser.Schema {}
			parser.File("./" + conf_schema.Get("dir").(string) + "/" + schema + ".yml")

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
