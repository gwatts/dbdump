package main

/*
TODO:
accept env var for cmd line arguments
make compatible with postgres
*/

import (
	"dbdump/dbdump"
	"flag"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const (
	defaultDbtype = "mysql"
	defaultHost   = "localhost"
	defaultPort   = 3306
)

var (
	addHeader       bool
	format          string
	dbtype          string
	envArgs         bool
	dbConfig        dbdump.DBConfig
	formatterConfig dbdump.FormatterConfig
)

func init() {
	flag.StringVar(&dbtype, "dbtype", defaultDbtype, "Database type.  Available databases: "+strings.Join(dbdump.DBNames(), ", "))
	flag.StringVar(&dbConfig.Host, "host", defaultHost, "Database hostname")
	flag.IntVar(&dbConfig.Port, "port", defaultPort, "Database port number")
	flag.StringVar(&dbConfig.Username, "user", "", "Database username")
	flag.StringVar(&dbConfig.Password, "password", "", "Database password")
	flag.StringVar(&dbConfig.Database, "database", "", "Database to connect to")
	flag.StringVar(&format, "format", "csv", "Format to dump into.  Available formatters: "+strings.Join(dbdump.FormatterNames(), ", "))
	flag.BoolVar(&formatterConfig.AddHeader, "header", true, "Enable to add column headers to the output (cvs)")
	flag.BoolVar(&envArgs, "env", false, "Enable to cause query parameters in the format $envname or ${envname} to be expanded from environment variables")
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s [flags] <query> [query-arg]...:\n", os.Args[0])
	flag.PrintDefaults()
}

func fail(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Println()
	os.Exit(1)
}

// expandArgs expands $envname or ${envname} from environment variables if the -env flag is set
func expandArgs(src []string, envArgs bool) []interface{} {
	result := make([]interface{}, len(src))
	for i, arg := range src {
		if envArgs {
			result[i] = os.ExpandEnv(arg)
		} else {
			result[i] = arg
		}
	}
	return result
}

func main() {
	flag.Usage = Usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		Usage()
		os.Exit(0)
	}

	if _, ok := dbdump.Formatters[format]; !ok {
		fail("Unknown formatter")
	}

	// connect and execute query
	newDb := dbdump.Databases[dbtype]
	if newDb == nil {
		fail("Invalid database type")
	}

	db, err := newDb(dbConfig)
	if err != nil {
		fail("Failed to connect to the database: %s", err)
	}

	qargs := expandArgs(args[1:], envArgs)

	stmt, err := db.Prepare(args[0])
	if err != nil {
		fail("Failed to prepare query: %v", err)
	}

	rows, err := stmt.Query(qargs...)

	if err != nil {
		fail("Query %q failed: %s", args[0], err)
	}

	cols, err := rows.Columns()
	if err != nil {
		fail("No results: %v", err)
	}

	w := dbdump.Formatters[format](formatterConfig, os.Stdout, cols)

	colcount := len(cols)
	rowvals := make([]interface{}, colcount)
	rowptrs := make([]interface{}, colcount)
	for i, _ := range rowvals {
		rowptrs[i] = &rowvals[i]
	}

	var i int
	for rows.Next() {
		err := rows.Scan(rowptrs...)
		if err != nil {
			fail("Failed to scan row %d: %s", i, err)
		}
		if err := w.Write(rowvals); err != nil {
			fail("Failed to write row %d: %s", i, err)
		}
		i++
	}
	w.Close()
}
