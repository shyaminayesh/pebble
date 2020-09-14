package schema

import (
	"os"
	"log"
	"fmt"
	"strings"
	"io/ioutil"
	"path/filepath"
	"gopkg.in/yaml.v3"
)

type (
	Schema struct {
		Name		string
		Path		string
		Structure	Structure
	}

	Structure struct {
		Table		Table		`yaml:"table"`
		Columns		[]Column	`yaml:"columns"`
		Keys		[]Key		`yaml:"keys"`
		Indexes		[]Index		`yaml:"indexes"`
	}

	Table struct {
		Engine		string	`yaml:"engine"`
		Charset		string	`yaml:"charset"`
		Collation	string	`yaml:"collation"`
	}

	Column struct {
		Name		string		`yaml:"name"`
		Type		string		`yaml:"type"`
		Primary		bool		`yaml:"primary"`
		Nullable	bool		`yaml:"nullable"`
		Increment	bool		`yaml:"increment"`
	}

	Key struct {
		Field		string		`yaml:"field"`
		Type		string		`yaml:"type"`
	}

	Index struct {
		Field		string		`yaml:"field"`
	}
)


func (schema *Schema) File(path string) {

	if _, err := os.Stat(path);
	os.IsNotExist(err) { log.Fatal(err) }

	/*
	We have a working schema file and now we can
	set known properties of the schema using the
	migration file
	*/
	schema.Path, _ = filepath.Abs(path)
	schema.Name = strings.Replace(filepath.Base(schema.Path), ".yml", "", 1)

	/*
	We need to read the file and append the table
	properties to the Schema struct to continue
	the migration file parse.
	*/
	buffer, err := ioutil.ReadFile(schema.Path)
	if err != nil { log.Fatal(err) }

	structure := Structure {}
	yaml.Unmarshal(buffer, &structure)
	schema.Structure = structure

}


func (schema *Schema) Statement() string {
	return fmt.Sprintf("CREATE TABLE `%s` (%s) ENGINE=%s DEFAULT CHARSET=%s DEFAULT COLLATE=%s", schema.Name, schema.ColumnStatement(), schema.Structure.Table.Engine, schema.Structure.Table.Charset, schema.Structure.Table.Collation)
}



/*
Generate sql statment lines for each column we have
in our migration answer file to append to the final
statement.
*/
func (schema *Schema) ColumnStatement() string {

	var statement string
	for _, column := range schema.Structure.Columns {

		// BASE
		sql := fmt.Sprintf("`%s` %s", column.Name, column.Type)

		// NULLABLE
		if column.Nullable == true { sql = sql + " NULL" }
		if column.Nullable == false { sql = sql + " NOT NULL" }

		// AUTO INCREMENT
		if column.Increment == true { sql = sql + " AUTO_INCREMENT" }

		// APPEND
		statement = statement + sql + ", "

	}

	/*
	Append the keys and the indexes statements into
	the main schema migration statement.
	*/
	is_keys, keys_stmnt := schema.KeysStatement()
	is_indexes, indexes_stmnt := schema.IndexesStatement()

	if is_keys && is_indexes {
		statement = statement + keys_stmnt
	} else if is_keys {
		statement = statement + strings.TrimSuffix(keys_stmnt, ",")
	}

	if is_indexes {
		statement = statement + indexes_stmnt
	}

	// RETURN
	return strings.TrimSuffix(statement, ", ")

}



func (schema *Schema) KeysStatement() (bool, string) {

	/*
	We have to check if there is keys present in the migration
	file and apply them to the final query.
	*/
	if len(schema.Structure.Keys) > 0 {

		var statement string
		for _, column := range schema.Structure.Keys {

			// STRING
			var sql string

			// PRIMARY
			if ( column.Type == "PRIMARY" ) {
				sql = fmt.Sprintf("  %s KEY (`%s`)", column.Type, column.Field)
			} else {
				sql = fmt.Sprintf("  %s KEY `%s` (`%s`)", column.Type, column.Field, column.Field)
			}

			statement = statement + sql + ",\n"
		}
		return true, strings.TrimSuffix(statement, "\n")

	} else {
		return false, ""
	}

}



func (schema *Schema) IndexesStatement() (bool, string) {

	/*
	We have to check if there is indexes available for this
	schema and build the query depending on the indexes.
	*/
	if len(schema.Structure.Indexes) > 0 {

		var statement string
		for _, column := range schema.Structure.Indexes {
			sql := fmt.Sprintf("INDEX `%s` (`%s`)", column.Field, column.Field)
			statement = statement + sql + ", "
		}
		return true, statement

	} else {
		return false, ""
	}

}