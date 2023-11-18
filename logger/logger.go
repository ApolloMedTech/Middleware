package logger

import (
	"github.com/ApolloMedTech/Middleware/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

func SetupLogger(cfg config.LogConfig) {
	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logLevel = logrus.InfoLevel // Default to InfoLevel if parsing fails
	}

	logrus.SetLevel(logLevel)

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp:       false,
		TimestampFormat:        "2006-01-02 15:04:05",
		DisableColors:          false,
		QuoteEmptyFields:       true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
		FullTimestamp:          false,
		// Customizing delimiters
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "caller",
		},
	})

	var logWriters []io.Writer
	logWriters = append(logWriters, &lumberjack.Logger{
		Filename:   cfg.LogPath,
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAgeDays,
	})

	if cfg.LogToStdout {
		logWriters = append(logWriters, os.Stdout)
	}

	multiWriter := io.MultiWriter(logWriters...)
	logrus.SetOutput(multiWriter)
}