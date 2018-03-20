package sqler

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/gohxs/sqler/sqler/sqlerdb"
	"github.com/gohxs/termu"

	// Drivers
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func New() *SQLer {
	return &SQLer{}
}

type SQLer struct {
	db   sqlerdb.DB
	term *termu.Term
}

type FileHistory struct {
	fileName string
	termu.History
}

func (fh *FileHistory) Load(fileName string) {
	fh.fileName = fileName
	data, err := ioutil.ReadFile(fh.fileName)
	if err != nil {
		return
	}
	parts := strings.Split(string(data), "\n")
	for _, p := range parts {
		fh.History.Append(p)
	}

}
func (fh *FileHistory) Append(in string) {
	fh.History.Append(in)

	f, err := os.OpenFile(fh.fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0600))
	if err != nil {
		return
	}
	f.WriteString(in + "\n")
	f.Close() // Slow?
}

func (s *SQLer) Start() {
	s.term = termu.New()
	s.term.SetPrompt("sqler> ")

	// Temp, sample commands for auto completer
	s.term.History.Append("select * from")
	s.term.History.Append(`.open cockroachdb postgres://root@localhost:2222?sslmode=disable`)
	s.term.History.Append(`.open cockroachdb postgres://root@localhost:26257`)
	s.term.History.Append(`.open postgres postgres://postgres@localhost?sslmode=disable`)
	s.term.History.Append(`.open sqlserver sqlserver://dev:1q2w3e@192.168.0.2`)

	fh := FileHistory{"", s.term.History}
	usr, _ := user.Current()
	dir := usr.HomeDir
	fh.Load(dir + "/.sqler_history")

	s.term.History = &fh

	ce := ComplEngine{term: s.term}
	s.term.AutoComplete = ce.AutoComplete
	s.term.Display = ce.Display
	ce.Suggest = s.sqlComplete()

	if len(os.Args) > 2 {
		drv := os.Args[1]
		dsn := os.Args[2]
		s.open(drv, dsn)

	}
	// Do something to history
	for {
		line, err := s.term.ReadLine()
		if err != nil {
			if err != termu.ErrEOF {
				log.Println("Err:", err)
			}
			return
		}
		s.runCmd(line) // Process command

		//fmt.Fprintln(s.term, "Echo:", line)
		// Process line here
	}
}

func (s *SQLer) runCmd(line string) {
	wr := &FmtWriter{s.term}
	args := parseargs(line)
	switch args[0] {
	case ".open": // Move operation to another area
		if len(args) < 3 {
			wr.Println("params < 3")
			return
		}
		driver := args[1]
		connStr := args[2]

		if err := s.open(driver, connStr); err != nil {
			wr.Println("Error connecting: ", err)
		}
		return

	case ".db": // Select database wrapper
		if len(args) != 2 {
			wr.Println("Requires 1 param")
			return
		}
		err := s.db.SetDatabase(args[1])
		if err != nil {
			wr.Println("E:", err)
		}
	case ".tables":
		if s.db == nil {
			wr.Println("Please connect to database to proceed")
			return
		}
		tables, err := s.db.ListTables()
		if err != nil {
			wr.Println("E:", err)
			return
		}
		dt := DataTable{}
		dt = append(dt, []interface{}{string("tables")})
		for _, t := range tables {
			dt = append(dt, []interface{}{t})
		}
		TableFprint(wr, dt)
	default:
		s.execSQL(line)
	}

}

func (s *SQLer) open(driver string, connStr string) error {
	w := &FmtWriter{s.term}

	w.Printf("Connecting to %s with '%s'\n", driver, connStr)
	var err error
	s.db, err = sqlerdb.New(driver, connStr)
	if err != nil {
		return err
	}

	v, _ := s.db.Version()
	w.Println(v)
	return nil

}
func (s *SQLer) execSQL(sqlStr string) {
	w := FmtWriter{s.term}

	if s.db == nil {
		w.Println("Please connect to database to proceed")
		return
	}
	result, err := s.db.Query(sqlStr)
	if err != nil {
		w.Println("Error:", err)
		return
	}
	defer result.Close()

	// Handle results
	resulti := 0
	for {
		var dataTable DataTable
		var vals []interface{}

		cols, _ := result.Columns()
		if len(cols) == 0 {
			w.Println("OK")
		}

		vals = make([]interface{}, len(cols))
		//colMaxWidth = make([]int, len(cols))

		// Init vals for reuse
		for i := range vals {
			var ns sql.NullString
			vals[i] = &ns
		}
		if len(cols) == 0 {
			result.Next()
			return
		}

		// Add headers to our memory table
		rowData := []interface{}{}
		for _, v := range cols {
			rowData = append(rowData, v)
		}
		dataTable = append(dataTable, rowData)

		// Load data
		for result.Next() {
			err := result.Scan(vals...)
			if err != nil {
				log.Println("Scan failed: ", err)
			}
			rowData := make([]interface{}, len(vals))
			for i, v := range vals {
				var val string
				if v.(*sql.NullString).Valid {
					val = v.(*sql.NullString).String
				} else {
					val = ""
				}
				rowData[i] = val
			}
			//fmt.Fprintln(wr, strings.Join(rowStr, " "))
			dataTable = append(dataTable, rowData)
		}
		TableFprint(w, dataTable)

		if !result.NextResultSet() {
			break
		}
		resulti++
	}
}

func (s *SQLer) ListDatabases() []string {
	if s.db == nil {
		fmt.Fprintln(s.term, "Not connected")
		return []string{}
	}

	ret, err := s.db.ListDatabases()
	if err != nil {
		fmt.Fprintln(s.term, "Err:", err)
		return []string{}
	}
	return ret
}

func (s *SQLer) ListTables() []string {
	if s.db == nil {
		fmt.Fprintln(s.term, "Not connected")
		return []string{}
	}

	ret, err := s.db.ListTables()
	if err != nil {
		fmt.Fprintln(s.term, "Err:", err)
		return []string{}
	}
	return ret
}

// completer
func (s *SQLer) sqlComplete() func(string) ([]string, string) {

	selectPrefix := citem("select",
		citem(
			func(in string) []string {
				if strings.TrimSpace(in) != "" {
					return []string{in}
				}
				return []string{"*"}
			}, citem("from", citem(s.ListTables, citem("where"))),
		),
	)

	updatePrefix := citem("update",
		citem(s.ListTables,
			citem("set"),
		),
	)
	deletePrefix := citem("delete",
		citem("from",
			citem(s.ListTables, citem("where")),
		),
	)
	insertPrefix := citem("insert",
		citem("into",
			citem(s.ListTables, citem("values")),
		),
	)
	createPrefix := citem("create",
		citem("database"),
		citem("table"),
	)
	dropPrefix := citem("drop",
		citem("database", citem(s.ListDatabases)),
		citem("table", citem(s.ListDatabases)),
	)

	//ROOT
	main := citem(
		citem("use", citem(s.ListDatabases)),

		// open
		citem(".open",
			citem("cockroachdb"),
			//citem(`postgres://root@localhost:2222?sslmode=disable`),
			//citem(`postgres://root@localhost:26257`),

			citem("postgres"),
			citem("mysql"),
			citem("sqllite"),
			citem("sqlserver",
				citem(`sqlserver://dev:1q2w3e@192.168.0.2`),
			),
		),
		// set database
		citem(".db", citem(s.ListDatabases)),
		citem(".tables"),

		selectPrefix,
		updatePrefix,
		deletePrefix,
		insertPrefix,
		createPrefix,
		dropPrefix,

		citem("explain",
			selectPrefix,
			updatePrefix,
			deletePrefix,
			insertPrefix,
			createPrefix,
			dropPrefix),
		citem("show", citem("databases"), citem("tables")),
		citem("begin"),
		citem("rollback"),
		citem("commit"),
	)

	//re := regexp.MustCompile("\\s+")
	//_ = re
	return func(line string) ([]string, string) {
		parts := parseargs(line)
		//parts := re.Split(line, -1)
		return main.search(parts)
	}
}
