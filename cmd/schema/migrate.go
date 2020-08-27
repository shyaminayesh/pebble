package schema

import (
	"fmt"
	"log"
	"strings"
	"io/ioutil"
	"database/sql"
	"pebble/utils/log"
	"pebble/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	parser "pebble/utils/parser/schema"
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

			/*
				We can delete the columns that are not present in the
				schema file safely before processing any other action.
			*/
			var columns []string
			for _, column := range structure.Columns {
				columns = append(columns, column.Name)
			}
			query := fmt.Sprintf("SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='%s' AND TABLE_SCHEMA='%s' AND COLUMN_NAME NOT IN ('%s')", schema, conf_connection.Get("name"), strings.Join(columns, "','"))
			rows, err := db.Query(query)
			if err != nil { log.Fatal(err) }
			defer rows.Close()

			for rows.Next() {
				var name string
				err := rows.Scan(&name)
				if err != nil { log.Fatal(err) }
				query := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", schema, name)
				db.Exec(query)
			}


			/*
				It's time to check each colum for changes and apply them
				to the live database
			*/
			var last_column string
			for _, column := range structure.Columns {

				type (
					Column struct {
						Name		sql.NullString
						Type		string
						Nullable	string
						Key			string
						Charset		sql.NullString
						Collation	sql.NullString
						Extra		string
					}
				)

				Result := Column {}
				query := fmt.Sprintf("SELECT `COLUMN_NAME`, `COLUMN_TYPE`, `IS_NULLABLE`, `COLUMN_KEY`, `CHARACTER_SET_NAME`, `COLLATION_NAME`, `EXTRA` FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE TABLE_NAME='%s' AND TABLE_SCHEMA='%s' AND COLUMN_NAME='%s'", schema, conf_connection.Get("name"), column.Name)
				db.QueryRow(query).Scan(&Result.Name, &Result.Type, &Result.Nullable, &Result.Key, &Result.Charset, &Result.Collation, &Result.Extra)

				// ADDITIONS
				if Result.Name.Valid == false {

					// BASE STATEMENT
					stmnt := fmt.Sprintf("ALTER TABLE `%s`.`%s`  ADD `%s` %s", conf_connection.Get("name"), schema, column.Name, column.Type)

					// COLLATION
					if len(column.Collation) > 0 { stmnt = fmt.Sprintf("%s COLLATE %s", stmnt, column.Collation) }

					// NULLABLE
					if column.Nullable == true { stmnt = fmt.Sprintf("%s NULL", stmnt) }
					if column.Nullable == false { stmnt = fmt.Sprintf("%s NOT NULL", stmnt) }

					// COLUMN ORDER
					stmnt = fmt.Sprintf("%s AFTER `%s`", stmnt, last_column)

					// EXECUTE
					db.Exec(stmnt)
				}

				// VALIDATE ( Type )
				if column.Type != Result.Type { db.Exec(fmt.Sprintf("ALTER TABLE `%s`.`%s` MODIFY COLUMN `%s` %s", conf_connection.Get("name"), schema, column.Name, column.Type)) }

				// VALIDATE ( Nullable )
				if Result.Nullable == "YES" { if column.Nullable == false { db.Exec(fmt.Sprintf("ALTER TABLE `%s`.`%s` MODIFY COLUMN `%s` %s NOT NULL", conf_connection.Get("name"), schema, column.Name, column.Type)) } }
				if Result.Nullable == "NO" { if column.Nullable == true { db.Exec(fmt.Sprintf("ALTER TABLE `%s`.`%s` MODIFY COLUMN `%s` %s NULL", conf_connection.Get("name"), schema, column.Name, column.Type)) } }

				/*
					Set last column we work on to help append columns with
					ordering in the next cycle
				*/
				last_column = Result.Name.String

			}

		}

	}

}