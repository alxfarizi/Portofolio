package golang_logging

import (
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestBasicLogging(t *testing.T) {
	logrus.New()
	logrus.SetLevel(logrus.DebugLevel)

	logrus.Debug("This is Debug")
	logrus.Info("This is Info")
	logrus.Warn("This is Warn")
	logrus.Error("This is Error")
}

func TestUserLogin(t *testing.T) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.WithFields(logrus.Fields{
		"user_id":    "12345",
		"email":      "user@example.com",
		"ip_address": "197.168.1.1",
		"action":     "login",
	}).Info("User Login")
}

func TestTextFormatter(t *testing.T) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
}

func TestLogEntry(t *testing.T) {
	logger := logrus.New()
	entry := logrus.NewEntry(logger)

	entry = entry.WithFields(logrus.Fields{
		"service": "auth",
		"version": "1.0.0",
	})
	entry.Info("User Login")
	entry.Warn("User Logout")
	entry.Error("User Login Failed")
}

func TestErrorLogging(t *testing.T) {
	logger := logrus.New()
	err := errors.New("Error")
	logger.WithError(err).Error("Error Logging")
}
