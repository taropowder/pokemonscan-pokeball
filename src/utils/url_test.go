package utils

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestGetUrlInfo(t *testing.T) {
	hash, code, length, err := GetUrlInfo("http://pokemon.taropowder.cn/123123123")
	if err == nil {
		log.Info(hash, code, length)
	} else {
		log.Warn(err)
	}
}
