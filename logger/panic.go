package logger

import (
	"fmt"
	"strings"
)

// Panic logs a message with panic level and appends a stack trace to the end of the message.
// It formats the message with a "[PANIC]" prefix and logs it using the panic level.
func (l *Log) Panic(message ...any) {
	var msg strings.Builder
	if len(message) > 0 {
		msg.WriteString("[PANIC] ")
		msg.WriteString(fmt.Sprintf("%v", message[0]))
		message[0] = msg.String()
	}
	l.newLog.Panic().Str("category", l.flag).Msg(fmt.Sprint(message...))
}