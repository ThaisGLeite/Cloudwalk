// logger/logger.go
package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Log is a global logger
var Log = logrus.New()

func init() {
	// Output to stdout instead of the default stderr and set the log level to Info.
	Log.Out = os.Stdout
	Log.SetLevel(logrus.InfoLevel)

	// Use JSON formatter
	Log.SetFormatter(&logrus.JSONFormatter{})
}
