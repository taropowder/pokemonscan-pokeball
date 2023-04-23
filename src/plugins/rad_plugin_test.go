package plugins

import (
	"testing"
)

func TestRad(t *testing.T) {
	p := RadPlugin{}
	err := p.Register(nil, ``)
	if err != nil {
		t.Error(err)
	}
	err = p.Run(0, `{"target":"https://x.x.x.x","allow_domains":"baidu.com,a.com,0.0.0.0/0"}`)
	if err != nil {
		t.Error(err)
	}
	//, vuls, err := p.GetResult(0)
	//if err != nil {
	//	t.Error(err)
	//}
	//if len(vuls.Vuls) != 0 {
	//	t.Error()
	//}
	//fmt.Println(vuls)
}
