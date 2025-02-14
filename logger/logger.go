package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type Log struct {
	newLog *zerolog.Logger
	flag   string
}

// New returns a new logger that will print to the console with the following format:
//
// [TIME] [LEVEL] [PREFIX] MESSAGE
//
// The time format is "2006/01/02 15:04:05".
func New(prefix string) *Log {

	zerolog.TimeFieldFormat = "2006/01/02 15:04:05"

	output := zerolog.ConsoleWriter{
		Out:           os.Stdout,
		TimeFormat:    "2006/01/02 15:04:05",
		PartsOrder:    []string{"time", "category", "message"},
		FieldsExclude: []string{"category"},
		FormatLevel: func(i any) string {
			return "[" + strings.ToUpper(i.(string)) + "]"
		},
		FormatTimestamp: func(i any) string {
			return "[" + i.(string) + "]"
		},
		FormatFieldName: func(i any) string {
			return ""
		},
		FormatFieldValue: func(i any) string {
			if category, ok := i.(string); ok && category != "" {
				return "[" + category + "]"
			}
			return ""
		},
		FormatMessage: func(i any) string {
			return i.(string)
		},
	}

	logger := zerolog.New(output).With().Timestamp().Logger()

	return &Log{
		newLog: &logger,
		flag:   prefix,
	}
}
