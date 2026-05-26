package logger

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm/logger"
)

var log = New("DB")

type CustomLogger struct {
	logger.Config
}

func (c *CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *c
	newLogger.LogLevel = level
	return &newLogger
}

func (c *CustomLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if c.LogLevel >= logger.Info {
		requestID := GetRequestID(ctx)
		if requestID != "" {
			fmt.Printf("[INFO] "+msg+" [%s]\n", append(data, requestID)...)
		} else {
			fmt.Printf("[INFO] "+msg+"\n", data...)
		}
	}
}

func (c *CustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if c.LogLevel >= logger.Warn {
		requestID := GetRequestID(ctx)
		if requestID != "" {
			fmt.Printf("[WARN] "+msg+" [%s]\n", append(data, requestID)...)
		} else {
			fmt.Printf("[WARN] "+msg+"\n", data...)
		}
	}
}

func (c *CustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if c.LogLevel >= logger.Error {
		requestID := GetRequestID(ctx)
		if requestID != "" {
			fmt.Printf("[ERROR] "+msg+" [%s]\n", append(data, requestID)...)
		} else {
			fmt.Printf("[ERROR] "+msg+"\n", data...)
		}
	}
}

func (c *CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if c.LogLevel <= 0 {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && c.LogLevel >= logger.Error:
		sql, _ := fc()
		log.ErrorWithoutTrace(ctx, fmt.Sprintf("%s %s", err, sql))
	case elapsed > c.SlowThreshold && c.SlowThreshold != 0 && c.LogLevel >= logger.Warn:
		sql, rows := fc()
		log.Log(ctx, fmt.Sprintf("SLOW SQL >= %v [%.3fms] [rows:%v] %s", c.SlowThreshold, float64(elapsed.Nanoseconds())/1e6, rows, sql))
	case c.LogLevel >= logger.Info:
		sql, rows := fc()
		log.Log(ctx, fmt.Sprintf("[%.3fms] [rows:%v] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql))
	}
}
