package main

import (
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func GetRandomUUID() (nuuid string, e error) {
	uid, e := uuid.NewRandom()
	if e != nil {
		return nuuid, e
	}
	nuuid = uid.String()
	return nuuid, e
}

func LogError(ce error) error {
	log.WithField("Timestamp", time.Now().Format(time.RFC1123)).Error(ce.Error())
	return ce
}
