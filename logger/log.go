package logger

import (
	"fmt"
	"strings"
)

// Log for logs a message with info level. It formats the message with an "[INFO]" prefix
// and logs it using the info level. The message is prepended with the category
// specified by the flag field.
// Note: If it used for logging DB, log level must be log info in CustomLogger struct.
func (l *Log) Log(message ...any) {
	var msg strings.Builder
	if len(message) > 0 {
		msg.WriteString("[INFO] ")
		msg.WriteString(fmt.Sprintf("%v", message[0]))
		message[0] = msg.String()
	}
	
	l.newLog.Info().Str("category", l.flag).Msg(fmt.Sprint(message...))
}
