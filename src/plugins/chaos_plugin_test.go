package plugins

import (
	"fmt"
	plugin_proto "pokemonscan-pokeball/src/proto/proto_struct/plugin"
	"sync"
	"testing"
)

func TestChaos(t *testing.T) {
	p := ChaosPlugin{Name: plugin_proto.ChaosPluginName, WorkingTasks: &sync.Map{}}
	err := p.Register(nil, ``)
	if err != nil {
		t.Error(err)
	}
	err = p.Run(0, `{"key":"","target":"jd.com"}`)
	if err != nil {
		t.Error(err)
	}
	info, _, err := p.GetResult(0)
	if err != nil {
		t.Error(err)
	}
	//if len(vuls.Vuls) != 0 {
	//	t.Error()
	//}
	//fmt.Println(vuls)
	fmt.Println(info)

}
