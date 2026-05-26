package logger

import (
	"context"
	"fmt"
	"strings"
)

// Panic logs a message with panic level, including the request ID from context appended at the end.
func (l *Log) Panic(ctx context.Context, message ...any) {
	requestID := GetRequestID(ctx)

	var msg strings.Builder
	if len(message) > 0 {
		msg.WriteString("[PANIC] ")
		msg.WriteString(fmt.Sprintf("%v", message[0]))
		if requestID != "" {
			msg.WriteString(fmt.Sprintf(" [%s]", requestID))
		}
		message[0] = msg.String()
	}
	l.newLog.Panic().Str("category", l.flag).Msg(fmt.Sprint(message...))
}
