package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.Logger{
	Out:   os.Stdout,
	Level: logrus.DebugLevel,
	Formatter: &logrus.TextFormatter{
		FullTimestamp: true,
	},
}
