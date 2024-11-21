package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	log.WithFields(log.Fields{
		"Reason": "Fail Not Exists",
	}).Error("Open Fail")

	log.Info("Success")
}
