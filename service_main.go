package svc

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func ServiceMain(name, description string, configs ...ServiceConfig) {
	svc, err := NewService(name, description, configs...)
	if err != nil {
		log.WithFields(log.Fields{
			"action": "main",
			"status": "error",
			"error":  err,
		}).Error("Failed to initialize service")
		os.Exit(-1)
	}

	if err := svc.Execute(); err != nil {
		log.WithFields(log.Fields{
			"action": "main",
			"status": "service_error",
			"error":  err,
		}).Error("Error executing service")
		os.Exit(-2)
	}
}
