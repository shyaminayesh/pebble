package schema

import (
	"fmt"
	"log"
	"database/sql"
	"pebble/config"
	"github.com/spf13/cobra"

	_ "github.com/go-sql-driver/mysql"
)

var Schema_Status_Command = &cobra.Command{
	Use: "status",
	Short: "View schema status",
	Run: schema_status,
}

func schema_status(cmd *cobra.Command, args []string) {

	type (
		Schema struct {
			Field			string
			Type			string
			Collation		sql.NullString
			Null			string
			Key				string
			Default			sql.NullString
			Extra			string
			Privileges		string
			Comment			string
		}
	)

	/*
		We have to initialize viper configuration
		for future usage in the below sections.
	*/
	var conf = config.Config()
	conf_connection := conf.Sub("connection")

	dialect := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", conf_connection.Get("user"), conf_connection.Get("password"), conf_connection.Get("host"), conf_connection.Get("port"), conf_connection.Get("name"))
	db, err := sql.Open(conf_connection.Get("driver").(string), dialect)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()


	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}


	schema := Schema {}
	err = db.QueryRow("SHOW FULL COLUMNS FROM address").Scan(&schema.Field, &schema.Type, &schema.Collation, &schema.Null, &schema.Key, &schema.Default, &schema.Extra, &schema.Privileges, &schema.Comment)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println( schema.Field )


}