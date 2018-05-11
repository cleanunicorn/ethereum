package web3

import (
	log "github.com/sirupsen/logrus"
)

func SetLogLevel(level log.Level) {
	log.SetLevel(level)
}
