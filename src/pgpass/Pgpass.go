package pgpass

import (
	"bufio"
	"os"
	"runtime"
	"strings"
)

type Pgpass struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

func GetPgpassFilepath() string {
	if path := os.Getenv("PGPASSFILE"); path != "" {
		return path
	}

	switch runtime.GOOS {
	case "windows":
		return os.Getenv("APPDATA") + "\\postgresql\\pgpass.conf"
	case "linux":
		return "~/.pgpass"
	}
	return ""
}

func LoadDefaultSettings() ([]Pgpass, error) {
	filepath := GetPgpassFilepath()

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	pgs, err := LoadSettings(file)
	file.Close()

	if err != nil {
		return nil, err
	}

	return pgs, err
}

func LoadSettings(file *os.File) ([]Pgpass, error) {
	pgs := make([]Pgpass, 0)

	s := bufio.NewScanner(file)
	for s.Scan() {
		pg := parseFileLine(s.Text())
		pgs = append(pgs, pg)
	}

	return pgs, nil
}

func parseFileLine(line string) Pgpass {
	var pg Pgpass

	cols := strings.Split(line, ":")

	// TODO: consider escape

	for i, value := range cols {
		value = strings.TrimSpace(value)
		switch i {
		case 0:
			pg.Host = value
			break
		case 1:
			pg.Port = value
			break
		case 2:
			pg.Database = value
			break
		case 3:
			pg.User = value
			break
		case 4:
			pg.Password = value
			break
		}
	}

	return pg
}

func (pg *Pgpass) Match(host string, port string, database string, user string) bool {
	if pg.Host == "*" && pg.Host != host {
		return false
	}
	if pg.Port == "*" && pg.Port != port {
		return false
	}
	if pg.Database == "*" && pg.Database != database {
		return false
	}
	if pg.User == "*" && pg.User != user {
		return false
	}
	return true
}
