package plugins

import (
	"fmt"
	plugin_proto "pokemonscan-pokeball/src/proto/proto_struct/plugin"
	"sync"
	"testing"
)

func TestCdnChecker(t *testing.T) {
	p := CdnCheckerPlugin{Name: plugin_proto.CdnCheckerPluginName, WorkingTasks: &sync.Map{}}
	err := p.Register(nil, ``)
	if err != nil {
		t.Error(err)
	}
	err = p.Run(0, `{
  "ip": "103.254.188.41",
  "checker_config": {
    "certificates": {
      "black_domain": ["cdn","chinanetcenter.com"],
      "white_domain": []
    }
  }
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
