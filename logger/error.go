package logger

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
)

// Error logs a message with error level, including the request ID from context,
// and appends a stack trace to the end of the message.
func (l *Log) Error(ctx context.Context, message ...any) {
	requestID := GetRequestID(ctx)

	var msg strings.Builder
	if len(message) > 0 {
		msg.WriteString("[ERROR] ")
		msg.WriteString(fmt.Sprintf("%v", message[0]))
		if requestID != "" {
			msg.WriteString(fmt.Sprintf(" [%s]", requestID))
		}
		message[0] = msg.String()

		message = append(message, "\n\n", string(debug.Stack()))
	}
	l.newLog.Error().Str("category", l.flag).Msg(fmt.Sprint(message...))
}

// ErrorWithoutTrace logs a message with error level without appending a stack trace,
// including the request ID from context appended at the end.
func (l *Log) ErrorWithoutTrace(ctx context.Context, message ...any) {
	requestID := GetRequestID(ctx)

	var msg strings.Builder
	if len(message) > 0 {
		msg.WriteString("[ERROR] ")
		msg.WriteString(fmt.Sprintf("%v", message[0]))
		if requestID != "" {
			msg.WriteString(fmt.Sprintf(" [%s]", requestID))
		}
		message[0] = msg.String()
	}
	l.newLog.Info().Str("category", l.flag).Msg(fmt.Sprint(message...))
}
