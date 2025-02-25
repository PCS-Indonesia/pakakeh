# Logger Package

## About
The Logger package provides custom logging utilities that can be integrated with the Gin framework and Gorm ORM. It allows for customizable log formats and log levels to suit your application's needs.

## Installation
To install the Logger package, simply import it into your `.go` files:

```go
import "github.com/PCS-Indonesia/pakakeh/logger"
```

## Usage
There's a different usages for logging gin, gorm (database) and also debugging.
Please refer to the `example` prefix files in the repository for detailed usage logging gin and gorm.

### Debugging Logger
To use the logger for debug, you need to define first with New function to set the prefix at the output log.
Prefix must be written with screaming kebab case ("CREATE-TICKET", "HOHO-HIHE", "HIHANG-HOHENG").

```go
package main

import (
	"github.com/PCS-Indonesia/pakakeh/logger"
)

func createProblems() {
    log := logger.New("PROBLEM")

    a := []int{1,2,3,4,5}
    for i, item := range a {
        log.Log("Problem - ", i, " is : ", a)
        // or u can simply with original go
        log.Log(fmt.Sprintf("Problem - %d is : %d", i, a))
    }
}

```

The instance of logger New can be initialized with dependency injection and usage by its method. Here's the example below,

```go
import (
	"github.com/PCS-Indonesia/pakakeh/logger"
)

type Service struct {
	repository *repository.Repository
	log        *logger.Log
}

func New(repo *repository.Repository, logger *logger.Log) *Service {
	return &Service{
		repository: repo,
		log:        logger.New("Service"),
	}
}

func (s *Service) deleteProblem() { 
    s.log.Log("hehe boii")

    // Your code here...
}
```

### Gin Logger
The output
To use the custom Gin logger, you need to set it up in your Gin engine:

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/PCS-Indonesia/pakakeh/logger"
)

func main() {
	r := gin.Default()

	// Set the custom Gin logger
	r.Use(gin.LoggerWithFormatter(logger.GinLogger))

    // See the example_main.go for the all gin logs that need to be used

	r.GET("/example", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	r.Run()
}
```

### Gorm Logger
The custom Gorm logger can be integrated as follows:

```go
package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/PCS-Indonesia/pakakeh/logger"
)

func main() {
    newLogger := &CustomLogger{
		Config: logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Error, // for display error only
			Colorful:      true,
			// your additionals config here
		},
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.NewCustomLogger(),
	})

	if err != nil {
		panic("failed to connect database")
	}
}
```
