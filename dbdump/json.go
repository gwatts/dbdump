package dbdump

import (
	"encoding/json"
	"io"
)

func init() {
	RegisterFormatter("json", NewJSONFormatter)
}

type JSONFormatter struct {
	w       io.Writer
	columns []string
	values  map[string]interface{}
	first   bool
}

func NewJSONFormatter(config FormatterConfig, w io.Writer, columns []string) Formatter {
	w.Write([]byte("["))
	return &JSONFormatter{
		w:       w,
		columns: columns,
		values:  make(map[string]interface{}, len(columns)),
		first:   true,
	}
}

func (jw *JSONFormatter) Write(values []interface{}) error {
	for i, colname := range jw.columns {
		switch v := values[i].(type) {
		case []byte:
			jw.values[colname] = string(v)
		default:
			jw.values[colname] = v
		}
	}
	entry, err := json.MarshalIndent(jw.values, "", "    ")
	if err != nil {
		return err
	}
	if jw.first {
		jw.first = false
	} else {
		jw.w.Write([]byte(","))
	}
	jw.w.Write(entry)
	return nil
}

func (jw *JSONFormatter) Close() error {
	_, err := jw.w.Write([]byte("]\n"))
	return err
}
