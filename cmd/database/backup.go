package database

import (
	"fmt"
	"log"
	"database/sql"
	"pebble/config"
	"pebble/utils/dumper"
	"github.com/spf13/cobra"
)

var Database_Backup_Command = &cobra.Command{
	Use: "backup",
	Short: "Backup database",
	Run: run_backup,
}


func run_backup(cmd *cobra.Command, args []string) {

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
	if err != nil { log.Fatal(err) }
	defer db.Close()

	/**
	* We have to Construct new instance of SQL dumper to
	* continue the backup process of the database using
	* export method.
	*/
	dumper, err := dumper.Construct(db)
	if err != nil { log.Fatal(err) }

	fmt.Println( dumper.Export() )

}