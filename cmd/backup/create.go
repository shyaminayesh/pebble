package backup

import (
	"os"
	"fmt"
	"log"
	"time"
	"database/sql"
	"pebble/config"
	"pebble/utils/dumper"
	"github.com/spf13/cobra"
)

var Create_Command = &cobra.Command{
	Use: "create",
	Short: "Create a new backup of the database",
	Run: run_create,
}


func run_create(cmd *cobra.Command, args []string) {

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

	dump, err := dumper.Export()
	if err != nil { log.Fatal(err) }

	/**
	* We have to create a file inside the backup directory
	* and store the dump text into that file to complete
	* the database backup.
	*/
	timestamp := time.Now().Format(fmt.Sprintf(Config.Backup.File.Timestamp))

	file, err := os.Create(fmt.Sprintf("./%s/%s.sql", Config.Backup.Directory, timestamp))
	if err != nil { log.Fatal(err) }
	defer file.Close()

	file.WriteString(fmt.Sprintf("%s\n", dump))

}