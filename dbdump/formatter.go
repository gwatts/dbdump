package dbdump

import "io"

var (
	Formatters = map[string]NewFormatter{}
)

type FormatterConfig struct {
	AddHeader bool
}

type Formatter interface {
	Write(values []interface{}) error
	Close() error
}

type NewFormatter func(config FormatterConfig, w io.Writer, columns []string) Formatter

func RegisterFormatter(name string, newFormatter NewFormatter) {
	Formatters[name] = newFormatter
}

func FormatterNames() (f []string) {
	for k, _ := range Formatters {
		f = append(f, k)
	}
	return f
}
