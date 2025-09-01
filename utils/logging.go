package utils

import "github.com/sirupsen/logrus"

func ErrorLog(message, path, serviceName string) {
	logrus.WithFields(logrus.Fields{
		"path":         path,
		"service_name": serviceName,
	}).Error(message)
}

func WarningLog(message, path, serviceName string) {
	logrus.WithFields(logrus.Fields{
		"path":         path,
		"service_name": serviceName,
	}).Warn(message)
}

func InfoLog(message, path, serviceName string) {
	logrus.WithFields(logrus.Fields{
		"path":         path,
		"service_name": serviceName,
	}).Info(message)
}
