package plugins

import (
	gonmap "github.com/lair-framework/go-nmap"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestParserRes(t *testing.T) {
	b, err := os.ReadFile("res.xml")
	if err != nil {
		log.Error("err", err)
	}
	nmapRes, err := gonmap.Parse(b)
	log.Info(nmapRes)
}
