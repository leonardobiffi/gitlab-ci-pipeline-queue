package queue

import (
	"os"

	"github.com/sirupsen/logrus"
)

func Logger() (log *logrus.Logger) {
	log = logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	logLevel := os.Getenv("QUEUE_PIPELINE_LOG_LEVEL")
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.SetLevel(logrus.InfoLevel)
		return log
	}

	log.SetLevel(level)
	return log
}
