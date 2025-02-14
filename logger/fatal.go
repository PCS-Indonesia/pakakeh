package logger

import (
	"fmt"
	"strings"
)

// Fatal logs a message with fatal level and appends a stack trace to the end of the message.
// It formats the message with an "[FATAL]" prefix and logs it using the fatal level.
func (l *Log) Fatal(message ...any) {
	var msg strings.Builder
	if len(message) > 0 {
		msg.WriteString("[FATAL] ")
		msg.WriteString(fmt.Sprintf("%v", message[0]))
		message[0] = msg.String()
	}
	l.newLog.Fatal().Str("category", l.flag).Msg(fmt.Sprint(message...))
}
