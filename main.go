package main

import (
	"github.com/Sirupsen/logrus"
)

var (
	logger logrus.FieldLogger = logrus.StandardLogger()
)

const (
	version = "0.1"
)

func main() {
	logger.WithFields(logrus.Fields{"version": version}).Info("EMS")
}
