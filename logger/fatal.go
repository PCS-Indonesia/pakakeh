package logger

import (
	"context"
	"fmt"
	"strings"
)

// Fatal logs a message with fatal level, including the request ID from context appended at the end.
func (l *Log) Fatal(ctx context.Context, message ...any) {
	requestID := GetRequestID(ctx)

	var msg strings.Builder
	if len(message) > 0 {
		msg.WriteString("[FATAL] ")
		msg.WriteString(fmt.Sprintf("%v", message[0]))
		if requestID != "" {
			msg.WriteString(fmt.Sprintf(" [%s]", requestID))
		}
		message[0] = msg.String()
	}
	l.newLog.Fatal().Str("category", l.flag).Msg(fmt.Sprint(message...))
}
