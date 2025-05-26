package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func init() {
	Log.Out = os.Stdout
	Log.SetLevel(logrus.WarnLevel)
}

func SetDebug(debug bool) {
	if debug {
		Log.SetLevel(logrus.DebugLevel)
		Log.Debug("Debug logging enabled")
	} else {
		Log.SetLevel(logrus.WarnLevel)
	}
}
