package dumper

import (
	"fmt"
	"log"
	"bytes"
	"database/sql"
	"text/template"
)

type (
	Dumper struct {
		db		*sql.DB
	}

	Data struct {
		Server	Server
		Tables	[]Table
	}

	Server struct {
		Database	string
		Version		string
	}

	Table struct {
		Name		string
		SQL			string
	}
)

const tmpl = `
-- Pebble (SQL Export)
--
-- Server version: {{ .Server.Version }}

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: {{ .Server.Database }}
--

-- --------------------------------------------------------

{{ range .Tables }}
--
-- Table structure for table: {{ .Name }}
--

DROP TABLE IF EXISTS {{ .Name }};
{{ .SQL }}

--
-- Dumping data for table {{ .Name }}
--


{{ end }}
`


func Construct(db *sql.DB) (*Dumper, error) {
	return &Dumper {
		db: db,
	}, nil
}


func (dumper *Dumper) Export() (string) {

	/**
	* Get information about the server instance to inject
	* into the export template.
	*/
	var database string
	query := fmt.Sprintf("SELECT DATABASE() AS name")
	dumper.db.QueryRow(query).Scan(&database)

	var version string
	query = fmt.Sprintf("SELECT VERSION() AS version")
	dumper.db.QueryRow(query).Scan(&version)


	/**
	* Append every property of the server information into
	* Server struct to later use in Data struct
	*/
	server := Server {
		Database: database,
		Version: version,
	}


	query = fmt.Sprintf("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA='%s'", server.Database)
	rows, err := dumper.db.Query(query)
	if err != nil { log.Fatal(err) }
	defer rows.Close()

	tables := []Table {}
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil { log.Fatal(err) }

		/**
		* Get the table create schema query for each table and
		* append it to the Table struct to push it to the text
		* template at the end.
		*/
		result := []string{"table", "ddl"}
		dumper.db.QueryRow(fmt.Sprintf("SHOW CREATE TABLE %s", table)).Scan(&result[0], &result[1])

		tables = append(tables, Table {
			Name: table,
			SQL: result[1],
		})
		// fmt.Println( result[1] )
	}


	/**
	* Create new text template base on the defined
	* template in the top of this file and append
	* any extra data to the text template before rendering
	*/
	tpl, err := template.New("dumper").Parse(tmpl)
	if err != nil { log.Fatal(err) }

	buffer := bytes.Buffer{}
	tpl.Execute(&buffer, Data {
		Server: server,
		Tables: tables,
	})


	return buffer.String()
}