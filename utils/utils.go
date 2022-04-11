package utils

import (
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func HandleError() {
	if err := recover(); err != nil {
		logrus.Errorln(err)
		debug.PrintStack()
	}
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
