package utils

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestCheckDomain(t *testing.T) {
	log.Info(CheckDomain("192.asd"))
}
