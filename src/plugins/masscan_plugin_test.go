package plugins

import (
	"fmt"
	"testing"
)

func TestMasscan(t *testing.T) {
	p := MasscanPlugin{}
	err := p.Register(nil, ``)
	if err != nil {
		t.Error(err)
	}
	err = p.Run(0, `{"command_args":"-p 80-100 1.1.1.1"}`)
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
