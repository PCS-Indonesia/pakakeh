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
		fmt.Printf("[INFO] "+msg+"\n", data...)
	}
}

func (c *CustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if c.LogLevel >= logger.Warn {
		fmt.Printf("[WARN] "+msg+"\n", data...)
	}
}

func (c *CustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if c.LogLevel >= logger.Error {
		fmt.Printf("[ERROR] "+msg+"\n", data...)
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
		log.Log(fmt.Sprintf("%s %s\n", err, sql))
	case elapsed > c.SlowThreshold && c.SlowThreshold != 0 && c.LogLevel >= logger.Warn:
		sql, rows := fc()
		log.Log(fmt.Sprintf("SLOW SQL >= %v [%.3fms] [rows:%v] %s\n", c.SlowThreshold, float64(elapsed.Nanoseconds())/1e6, rows, sql))
	case c.LogLevel >= logger.Info:
		sql, rows := fc()
		log.Log(fmt.Sprintf("[%.3fms] [rows:%v] %s\n", float64(elapsed.Nanoseconds())/1e6, rows, sql))
	}
}
