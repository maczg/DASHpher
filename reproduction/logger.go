package reproduction

import (
	"github.com/sirupsen/logrus"
)

var logger = logrus.Logger{
	Formatter: &logrus.TextFormatter{DisableColors: false, TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true},
}
