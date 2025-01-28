package logger

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Exampledatabase() {
	var err error

	newLogger := &CustomLogger{
		Config: logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Error, // for display error only
			Colorful:      true,
			// your config here
		},
	}

	DB, err = gorm.Open(mysql.Open("root:password@tcp(127.0.0.1:3306)/test"), &gorm.Config{
		Logger: newLogger,
		// your config here
	})

	if err != nil {
		log.Fatal("Failed to connect to database")
	}
}
