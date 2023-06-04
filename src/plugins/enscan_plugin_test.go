package plugins

import (
	"fmt"
	plugin_proto "pokemonscan-pokeball/src/proto/proto_struct/plugin"
	"testing"
)

func TestEnscanPlugin(t *testing.T) {
	p := ENScanPlugin{Name: plugin_proto.ENScanPluginName}
	err := p.Register(nil, ``)
	if err != nil {
		t.Error(err)
	}
	err = p.Run(0, `{
            "target": "郑州商学院",
            "type": "aqc,tyc"
          }`)
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
