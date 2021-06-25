package log

// Definition of log system. Ideally we would like to contain this application
// in a docker container and contain logstash, to store it in a ElasticSearch database for example.

import (
	"fmt"
	"os"
	"strings"

	"github.com/nickmonks/microservices-go/src/api/config"
	"github.com/sirupsen/logrus"
)

var (
	Log *logrus.Logger
)

func init() {
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		level = logrus.DebugLevel
	}

	Log = &logrus.Logger{
		Level: level,
		Out:   os.Stdout,
	}

	// In production we want to have a json formater, but for development we just want
	// stdout
	if config.IsProduction() {
		Log.Formatter = &logrus.JSONFormatter{}
	} else {
		Log.Formatter = &logrus.TextFormatter{}
	}
}

func Info(msg string, tags ...string) {
	if Log.Level < logrus.InfoLevel {
		return
	}

	// Create personalized logrus fields by parsing the string of
	// key/value pairs.
	Log.WithFields(ParseFields(tags...)).Info(msg)
}

func Error(msg string, err error, tags ...string) {
	if Log.Level < logrus.InfoLevel {
		return
	}

	// Create personalized logrus fields by parsing the string of
	// key/value pairs.
	msg = fmt.Sprintf("%s -- ERROR -- %v", msg, err)
	Log.WithFields(ParseFields(tags...)).Error(msg)
}

func Debug(msg string, tags ...string) {
	if Log.Level < logrus.InfoLevel {
		return
	}

	// Create personalized logrus fields by parsing the string of
	// key/value pairs.
	Log.WithFields(ParseFields(tags...)).Debug(msg)
}

// function use to transform string tags to the logrus fields
func ParseFields(tags ...string) logrus.Fields {
	// logrus.Fields is a map, we allocate the size of the map to be
	// the length of tags
	result := make(logrus.Fields, len(tags))

	for _, tag := range tags {
		// from each tag, just extract the first word for KEY, and second word
		// for VALUE
		els := strings.Split(tag, ":")
		result[strings.TrimSpace(els[0])] = strings.TrimSpace(els[1])
	}

	return result

}
