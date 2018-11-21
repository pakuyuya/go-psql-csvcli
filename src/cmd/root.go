package cmd

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	_ "github.com/lib/pq"
)

type Option struct {
	Help       bool
	Host       string
	Port       int
	Database   string
	User       string
	Pass       string
	SQL        string
	Quote      string
	Sepalate   string
	Escape     bool
	EscapeType string
	Nullas     string
	QuoteNull  bool
	RequireSsl bool
}

var opt Option

func init() {
	rootCmd.PersistentFlags().BoolVarP(&opt.Help, "help", "", false, "help for "+rootCmd.Name())
	rootCmd.PersistentFlags().StringVarP(&opt.Host, "host", "h", "localhost", "host name to connect")
	rootCmd.PersistentFlags().IntVarP(&opt.Port, "port", "p", 5432, "database identifier to connect")
	rootCmd.PersistentFlags().StringVarP(&opt.Database, "database", "d", "postgres", "database identifier to connect")
	rootCmd.PersistentFlags().StringVarP(&opt.User, "user", "u", "postgres", "user name using authentication")
	rootCmd.PersistentFlags().StringVarP(&opt.Pass, "pass", "w", "", "password using authentication")
	rootCmd.PersistentFlags().StringVar(&opt.SQL, "sql", "", "sql to execute.")
	rootCmd.PersistentFlags().StringVar(&opt.Quote, "quote", "\"", "quote csv column")
	rootCmd.PersistentFlags().StringVarP(&opt.Sepalate, "sepalate", "s", ",", "quote csv column")
	rootCmd.PersistentFlags().BoolVarP(&opt.Escape, "escape", "e", true, "escape spacial characters if use quots")
	rootCmd.PersistentFlags().StringVar(&opt.EscapeType, "escapetype", "cascade", "escape method. 'cascade' or 'backslash'")
	rootCmd.PersistentFlags().StringVar(&opt.Nullas, "nullas", "", "output of null as csv data")
	rootCmd.PersistentFlags().BoolVar(&opt.QuoteNull, "quotenull", true, "put quotes on null column")
	rootCmd.PersistentFlags().BoolVar(&opt.RequireSsl, "requiressl", false, "use ssl forced")
}

var rootCmd = &cobra.Command{
	Use:   "pq-csv-console",
	Short: "connect to postgresql database and pull csv.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := adjustArg(); err != nil {
			fmt.Println(err.Error())
			return
		}

		if err := readAndWriteCSV(); err != nil {
			fmt.Println(err.Error())
			return
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func adjustArg() error {
	stdin := bufio.NewScanner(os.Stdin)
	if piped := stdin.Text(); piped != "" {
		opt.SQL = piped
	}

	return nil
}

func readAndWriteCSV() error {

	sslmode := "disable"
	if opt.RequireSsl {
		sslmode = "require"
	}
	dbSourceName := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		opt.Host, opt.Port, opt.User, opt.Pass, opt.Database, sslmode)

	db, err := sql.Open("postgres", dbSourceName)
	if err != nil {
		return err
	}

	rows, err := db.Query(opt.SQL)
	if err != nil {
		return err
	}

	columns, _ := rows.Columns()
	for rows.Next() {
		row := make([]interface{}, len(columns))
		for idx := range columns {
			row[idx] = new(MetalScanner)
		}

		err := rows.Scan(row...)
		if err != nil {
			return err
		}
		writeCSV(row, &columns)
	}

	return nil
}

func writeCSV(row []interface{}, columns *[]string) {
	bytes := make([]byte, 0)
	l := len(*columns)

	for idx := range *columns {
		var scanner = row[idx].(*MetalScanner)

		v, notnull := scanner.value.(string)

		if notnull {
			v = escapeValue(v)
			v = opt.Quote + v + opt.Quote
		} else {
			v = opt.Nullas
			if opt.QuoteNull {
				v = opt.Quote + v + opt.Quote
			}
		}

		bytes = append(bytes, v...)
		if idx < l-1 {
			bytes = append(bytes, opt.Sepalate...)
		}
	}

	fmt.Println(string(bytes))
}

func escapeValue(v string) string {
	if !opt.Escape {
		return v
	}
	switch opt.EscapeType {
	case "cascade":
		v = strings.Replace(v, "\"", "\"\"", -1)
		break
	case "backslash":
		v = strings.Replace(v, "\\", "\\\\", -1)
		v = strings.Replace(v, "\"", "\\\"", -1)
		break
	}

	return v
}

type MetalScanner struct {
	valid bool
	value interface{}
}

func (scanner *MetalScanner) getBytes(src interface{}) []byte {
	if a, ok := src.([]uint8); ok {
		return a
	}
	return nil
}
func (scanner *MetalScanner) Scan(src interface{}) error {
	switch src.(type) {
	case int64:
		if value, ok := src.(int64); ok {
			scanner.value = strconv.FormatInt(value, 10)
			scanner.valid = true
		}
	case float64:
		if value, ok := src.(float64); ok {
			scanner.value = fmt.Sprintf("%f", value)
			scanner.valid = true
		}
	case bool:
		if value, ok := src.(bool); ok {
			scanner.value = strconv.FormatBool(value)
			scanner.valid = true
		}
	case string:
		//value := scanner.getBytes(src)
		scanner.value = src.(string)
		scanner.valid = true
	case []byte:
		value := scanner.getBytes(src)
		scanner.value = string(value)
		scanner.valid = true
	case time.Time:
		if value, ok := src.(time.Time); ok {
			scanner.value = value.Format("2006-01-02 15:04:05")
			scanner.valid = true
		}
	case nil:
		scanner.value = nil
		scanner.valid = true
	}
	return nil
}
