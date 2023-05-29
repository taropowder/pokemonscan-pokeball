package utils

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestGetUrlInfo(t *testing.T) {
	hash, code, title, length, err := GetUrlInfo("https://zhoushan.dujia.qunar.com")
	if err == nil {
		log.Info(hash, title, code, length)
	} else {
		log.Warn(err)
	}
	hash, code, title, length, err = GetUrlInfo("https://zhejiang.dujia.qunar.com")
	if err == nil {
		log.Info(hash, title, code, length)
	} else {
		log.Warn(err)
	}
}
