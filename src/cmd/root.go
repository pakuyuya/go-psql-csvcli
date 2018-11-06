package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Option struct {
	Host       string
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
}

var opt Option

func init() {
	rootCmd.PersistentFlags().StringVarP(&opt.Host, "host", "h", "", "host name to connect")
	rootCmd.PersistentFlags().StringVarP(&opt.Database, "databse", "d", "", "database identifier to connect")
	rootCmd.PersistentFlags().StringVarP(&opt.User, "user", "u", "", "user name using authentication")
	rootCmd.PersistentFlags().StringVarP(&opt.Pass, "pass", "p", "", "password using authentication")
	rootCmd.PersistentFlags().StringVarP(&opt.SQL, "sql", "q", "", "sql to execute.")
	rootCmd.PersistentFlags().StringVarP(&opt.Quote, "quote", "", "\"", "quote csv column")
	rootCmd.PersistentFlags().StringVarP(&opt.Sepalate, "sepalate", "s", ",", "quote csv column")
	rootCmd.PersistentFlags().BoolVarP(&opt.Escape, "escape", "e", true, "escape spacial characters if use quots")
	rootCmd.PersistentFlags().StringVarP(&opt.EscapeType, "escapetype", "", "cascade", "escape method. 'cascade' or 'backslash'")
	rootCmd.PersistentFlags().StringVarP(&opt.Nullas, "nullas", "", "", "output of null as csv data")
	rootCmd.PersistentFlags().BoolVarP(&opt.QuoteNull, "nullas", "", true, "put quotes on null column")
}

var rootCmd = &cobra.Command{
	Use:   "pq-csv-console",
	Short: "connect to postgresql database and pull csv.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := adjustArg(); err != nil {
			fmt.Println(err.Error())
			return
		}

		if err := argCheck(); err != nil {
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
	// TODO:
	return nil
}

func argCheck() error {
	// TODO:
	return nil
}

func readAndWriteCSV() error {
	// TODO:
	return nil
}
