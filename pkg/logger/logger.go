package logger

import (
	"github.com/gryd-database/platform-poc/configuration"
	"os"

	"github.com/sirupsen/logrus"
)

var filepath = "./pkg/logger/logrus.log"

var logger_instance = logrus.New()

type loggerService struct {
	logger *logrus.Logger
}

func Init(configService *configuration.Config) (*logrus.Logger, error) {

	// setting the format of the logs to be a JSON one
	logger_instance.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})

	// getting the log level set in the configuration file
	logLevel, err := logrus.ParseLevel(configService.Logger.LogLevel)
	// If the log level in conf file can't be parsed, log level should be the default info level
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	// setting the log level
	logger_instance.SetLevel(logLevel)

	if configService.Logger.LogEnv == "local" { // If we want to throw logs into a local file

		logger_instance.SetOutput(os.Stdout)
		// setting it to a file writer
		file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

		if err == nil {
			logger_instance.Out = file
		} else {
			logger_instance.Info("Failed to log to file, using default stderr")
		}
	}
	return logger_instance, nil
}
