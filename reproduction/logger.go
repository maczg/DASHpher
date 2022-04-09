package reproduction

import (
	"github.com/sirupsen/logrus"
	"os"
)

var logger = logrus.Logger{
	Out:       os.Stderr,
	Formatter: &logrus.TextFormatter{DisableColors: false, TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true},
	Level:     logrus.InfoLevel,
}
