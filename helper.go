package arest

import log "github.com/sirupsen/logrus"

var IsDebug bool = false

func Debug(msg string, args ...interface{}) {
	if IsDebug {
		log.Debugf(msg, args...)
	}
}
