package logger

import (
	"context"
	"fmt"
	"strings"
)

// Log logs a message with info level, including the request ID from context.
// It formats the message with an "[INFO]" prefix and the request ID appended at the end.
func (l *Log) Log(ctx context.Context, message ...any) {
	requestID := GetRequestID(ctx)

	var msg strings.Builder
	if len(message) > 0 {
		msg.WriteString("[INFO] ")
		msg.WriteString(fmt.Sprintf("%v", message[0]))
		if requestID != "" {
			msg.WriteString(fmt.Sprintf(" [%s]", requestID))
		}
		message[0] = msg.String()
	}

	l.newLog.Info().Str("category", l.flag).Msg(fmt.Sprint(message...))
}
