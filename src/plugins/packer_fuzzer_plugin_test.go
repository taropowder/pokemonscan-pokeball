package plugins

import (
	"fmt"
	plugin_proto "pokemonscan-pokeball/src/proto/proto_struct/plugin"
	"testing"
)

func TestPackerFuzzer(t *testing.T) {
	p := PackerFuzzerPlugin{Name: plugin_proto.PackerFuzzerPluginsPluginName}
	err := p.Register(nil, ``)
	if err != nil {
		t.Error(err)
	}
	//err = p.Run(123, `{"target":"https://pokemon.taropowder.cn"}`)
	//if err != nil {
	//	t.Error(err)
	//}

	info, _, err := p.GetResult(123)
	if err != nil {
		t.Error(err)
	}
	//if len(vuls.Vuls) != 0 {
	//	t.Error()
	//}
	//fmt.Println(vuls)
	fmt.Println(info)

}
