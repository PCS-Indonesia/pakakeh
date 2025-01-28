package logger

import (
	"fmt"
	"runtime/debug"
	"strings"
)

// Error logs message with error level and appends a stack trace to the end of the message.
func (l *Log) Error(message ...any) {
	var msg strings.Builder
	if len(message) > 0 {
		msg.WriteString("[ERROR] ")
		msg.WriteString(fmt.Sprintf("%v", message[0]))
		message[0] = msg.String()

		message = append(message, "\n\n", string(debug.Stack()))
	}
	l.newLog.Error().Str("category", l.flag).Msg(fmt.Sprint(message...))
}


// ErrorWithoutTrace logs a message with error level without appending a stack trace.
// It formats the message with an "[ERROR]" prefix and logs it using the info level.
func (l *Log) ErrorWithoutTrace(message ...any) {
	var msg strings.Builder
	if len(message) > 0 {
		msg.WriteString("[ERROR] ")
		msg.WriteString(fmt.Sprintf("%v", message[0]))
		message[0] = msg.String()
	}
	l.newLog.Info().Str("category", l.flag).Msg(fmt.Sprint(message...))
}
