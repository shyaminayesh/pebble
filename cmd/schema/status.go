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

var conf = config.Config()

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
	config_connection := conf.Sub("connection")

	conuri := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", config_connection.Get("user"), config_connection.Get("password"), config_connection.Get("host"), config_connection.Get("port"), config_connection.Get("name"))
	db, err := sql.Open(config_connection.Get("driver").(string), conuri)
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