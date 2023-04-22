package plugins

import (
	"fmt"
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

func TestNmap(t *testing.T) {
	p := NmapPlugin{}
	err := p.Register(nil, ``)
	if err != nil {
		t.Error(err)
	}
	err = p.Run(0, `{"command_args":"-p 80 1.1.1.1"}`)
	if err != nil {
		t.Error(err)
	}
	info, vuls, err := p.GetResult(0)
	if err != nil {
		t.Error(err)
	}
	if len(vuls.Vuls) != 0 {
		t.Error()
	}
	fmt.Println(vuls)
	fmt.Println(info)

}
