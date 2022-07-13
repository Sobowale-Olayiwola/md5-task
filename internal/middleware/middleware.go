package middleware

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerToFile() gin.HandlerFunc {

	logFileName := os.Getenv("LOG_FILE_NAME")
	logFilePath := os.Getenv("LOG_FILE_PATH")

	fileName := path.Join(logFilePath, logFileName)

	//Write to file
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("err", err.Error())
	}

	// instantiation
	logger := logrus.New()

	//Set output
	logger.Out = src

	//Set log level
	logger.SetLevel(logrus.DebugLevel)

	//Format log
	logger.SetFormatter(&logrus.TextFormatter{})

	return func(c *gin.Context) {
		//Start time
		startTime := time.Now()

		//Process request
		c.Next()

		//End time
		endTime := time.Now()

		//Execution time
		latencyTime := endTime.Sub(startTime)

		//Request method
		reqMethod := c.Request.Method

		//Request routing
		reqUri := c.Request.RequestURI

		// status code
		statusCode := c.Writer.Status()

		// request IP
		clientIP := c.ClientIP()

		//Log format
		logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}
