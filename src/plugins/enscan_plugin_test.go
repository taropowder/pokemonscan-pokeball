package plugins

import (
	"encoding/base64"
	"fmt"
	"os"
	plugin_proto "pokemonscan-pokeball/src/proto/proto_struct/plugin"
	"testing"
)

func TestEnscanPlugin(t *testing.T) {
	p := ENScanPlugin{Name: plugin_proto.ENScanPluginName}
	err := p.Register(nil, ``)
	if err != nil {
		t.Error(err)
	}

	data, err := os.ReadFile("/tmp/poke_test/config.yaml")
	if err != nil {
		fmt.Println("读取文件错误:", err)
		return
	}

	base64Data := base64.StdEncoding.EncodeToString(data)

	err = p.Run(0, fmt.Sprintf(`{
            "target": "郑州商学院",
            "type": "aqc",
           "enscan_config_file": "%s"
          }`, base64Data))
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
