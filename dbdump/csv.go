package dbdump

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

func init() {
	RegisterFormatter("csv", NewCSVFormatter)
}

type CSVFormatter struct {
	w      *csv.Writer
	rowstr []string
}

func NewCSVFormatter(config FormatterConfig, w io.Writer, columns []string) Formatter {
	cw := &CSVFormatter{w: csv.NewWriter(w)}
	if config.AddHeader {
		cw.w.Write(columns)
	}
	cw.rowstr = make([]string, len(columns))
	return cw
}

func (cw *CSVFormatter) Write(values []interface{}) error {
	for j, col := range values {
		switch v := col.(type) {
		case []byte:
			cw.rowstr[j] = string(v)
		case string:
			cw.rowstr[j] = v
		case int64:
			cw.rowstr[j] = strconv.FormatInt(v, 10)
		case float64:
			cw.rowstr[j] = strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			cw.rowstr[j] = strconv.FormatBool(v)
		case time.Time:
			cw.rowstr[j] = v.String()
		case nil:
			cw.rowstr[j] = "NULL"
		default:
			return fmt.Errorf("Unknown type %T returned from database driver", col)
		}

	}
	if err := cw.w.Write(cw.rowstr); err != nil {
		return err
	}
	return nil
}

func (cw *CSVFormatter) Close() error {
	cw.w.Flush()
	return nil
}
