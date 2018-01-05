package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"hexasoftware/lib/cli-table"
	"io"
	"log"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/quick"
	"github.com/alecthomas/chroma/styles"
	"github.com/gohxs/termu"
	//"github.com/gohxs/readline"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	cmdH := cmdHandler{}

	cPrompt := "SQLi> "
	term := termu.New() // open
	term.SetPrompt(cPrompt)
	term.Display = display

	cmdH.wr = term

	//cmds := []string{}

	for {
		cmds := []string{}
		line, err := term.ReadLine()
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		cmds = append(cmds, line)
		if !strings.HasSuffix(line, ";") {
			n := len(cPrompt) - 3
			term.SetPrompt(strings.Repeat(" ", n) + ">> ")
			continue
		}
		cmd := strings.Join(cmds, "\n")
		cmds = cmds[:0]
		//rl.SetPrompt("SQLIr> ")
		term.SetPrompt(cPrompt)

		cmdH.Cmd(cmd)
	}

}

// Highlighter test
func display(input string) string {
	buf := bytes.NewBuffer([]byte{})
	err := quick.Highlight(buf, input, "postgres", "terminal16", "monokaim")
	//err := quick.Highlight(buf, input, "bash", "terminal16m", "monokaim")
	if err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

type cmdHandler struct {
	db *sql.DB
	wr io.Writer
}

// Cmd Sql command handler
func (h *cmdHandler) Cmd(line string) {

	wr := h.wr
	// Tokenize command
	re := regexp.MustCompile("[[:space:]]+")
	// Cool trick
	strPart := strings.Split(line, "\"")
	var args []string
	for i, v := range strPart {
		if i&0x1 == 0 {
			lparts := re.Split(strings.Trim(v, " "), -1)
			args = append(args, lparts...) // Add all parts if not inside '"'
		} else {
			args = append(args, v)
		}
	}

	if args[0] == ".open" {
		if len(args) < 3 {
			log.Println("Requires at least 2 arguments")
			return
		}
		var err error
		driver := args[1]
		connStr := args[2]

		if driver == "cockroachdb" { // alias
			driver = "postgres"
		}

		fmt.Fprintln(wr, "Connecting to "+driver+" '"+connStr+"'")
		h.db, err = sql.Open(driver, connStr)
		if err != nil {
			fmt.Fprintln(wr, "Error connecting: ", err)
			return
		}
		return
	}

	if h.db == nil {
		fmt.Fprintln(wr, "Please connect to database to proceed")
		return
	}

	result, err := h.db.Query(line)
	if err != nil {
		fmt.Fprintln(wr, "Error:", err)
		return
	}
	defer result.Close()

	// Read to array of strings first
	resulti := 0
	for result.NextResultSet() {
		var dataTable cliTable.DataTable

		var vals []interface{}
		cols, _ := result.Columns()
		if len(cols) == 0 {
			fmt.Fprintln(wr, "OK")
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
		cliTable.Fprint(wr, dataTable)

		resulti++
	}

	// PreBuild format string with the column paddings

}

// Monokai style.
var Monokai = styles.Register(chroma.MustNewStyle("monokaim", chroma.StyleEntries{
	chroma.Text:                "#f8f8f2",
	chroma.Error:               "#960050 bg:#1e0010",
	chroma.Comment:             "#75715e",
	chroma.Keyword:             "#66d9ef",
	chroma.KeywordNamespace:    "#f92672",
	chroma.Operator:            "#f92672",
	chroma.Punctuation:         "#f8f8f2",
	chroma.Name:                "#f8f8f2",
	chroma.NameBuiltin:         "#f844d0",
	chroma.NameAttribute:       "#a6e22e",
	chroma.NameClass:           "#a6e22e",
	chroma.NameConstant:        "#66d9ef",
	chroma.NameDecorator:       "#a6e22e",
	chroma.NameException:       "#a6e22e",
	chroma.NameFunction:        "#a6e22e",
	chroma.NameOther:           "#a6e22e",
	chroma.NameTag:             "#f92672",
	chroma.LiteralNumber:       "#ae81ff",
	chroma.Literal:             "#ae81ff",
	chroma.LiteralDate:         "#e6db74",
	chroma.LiteralString:       "#e6db74",
	chroma.LiteralStringEscape: "#ae81ff",
	chroma.GenericDeleted:      "#f92672",
	chroma.GenericEmph:         "italic",
	chroma.GenericInserted:     "#a6e22e",
	chroma.GenericStrong:       "bold",
	chroma.GenericSubheading:   "#75715e",
}))
