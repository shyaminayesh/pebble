package schema

import (
	"fmt"
	"log"
	"strings"
	"io/ioutil"
	"database/sql"
	"pebble/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		migration directory to find
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
		fmt.Println("[DROP]: " + table)
		db.Exec("DROP TABLE " + table)
	}

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
			Here we can handle table not exists state and we need to
			create the table accoding to the schema file provided.
		*/
		if count == 0 {

			var columns_stmnt string
			var columns_length = len(structure.Columns) - 1
			for index, column := range structure.Columns {

				// BASE STATEMENT
				columns_stmnt = columns_stmnt + column.Name + " " + column.Type

				// COLLATION
				if len(column.Collation) > 0 { columns_stmnt = columns_stmnt + " COLLATE " + column.Collation }

				// NULLABLE
				if column.Nullable == true { columns_stmnt = columns_stmnt + " NULL" }
				if column.Nullable == false { columns_stmnt = columns_stmnt + " NOT NULL" }

				// AUTO INCREMENT
				if column.Primary == true {
					if column.Increment == true { columns_stmnt = columns_stmnt + " PRIMARY KEY AUTO_INCREMENT" }
					if column.Increment == false { columns_stmnt = columns_stmnt + " PRIMARY KEY" }
				}

				// APPEND COMMA
				if index != columns_length { columns_stmnt = columns_stmnt + "," }

			}
			query := fmt.Sprintf("CREATE TABLE %s (%s) ENGINE=%s DEFAULT CHARSET=%s COLLATE=%s", schema, columns_stmnt, structure.Table.Engine, structure.Table.Charset, structure.Table.Collation)
			fmt.Println( query )
			db.Exec(query)

		}

		/*
			Here we need to handle table exisits and schema need to
			be validated state.
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
			query := fmt.Sprintf("SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='%s' AND TABLE_SCHEMA='%s' AND COLUMN_NAME NOT IN ('%s');", schema, conf_connection.Get("name"), strings.Join(columns, "','"))
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

		}

		// for rows.Next() {
		// 	var name string
		// 	err := rows.Scan(&name)
		// 	if err != nil { log.Fatal(err) }
		// 	fmt.Println("[DROP]: " + name)
		// 	// db.Exec("DROP TABLE " + table)
		// }

	}
	// "SHOW FULL COLUMNS FROM users WHERE Field NOT IN ('id')"

}