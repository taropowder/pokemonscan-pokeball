package plugins

import (
	"fmt"
	"testing"
)

func TestNuclei(t *testing.T) {
	p := NucleiPlugin{}
	err := p.Register(nil, ``)
	if err != nil {
		t.Error(err)
	}
	err = p.Run(0, `{"command_args":"-t baidu.cn -headless -severity critical"}`)
	if err != nil {
		t.Error(err)
	}
	_, vuls, err := p.GetResult(0)
	if err != nil {
		t.Error(err)
	}
	if len(vuls.Vuls) != 0 {
		t.Error()
	}
	fmt.Println(vuls)

}
