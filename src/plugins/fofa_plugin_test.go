package plugins

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
	"testing"
)

func TestFofa(t *testing.T) {
	p := FofaPlugin{}
	err := p.Register(nil, `{
  "email": "",
  "key": ""
}`)
	if err != nil {
		t.Error(err)
	}
	err = p.Run(0, `{
  "query": "domain=\"\"",
  "type": "query"
}`)
	if err != nil {
		t.Error(err)
	}

}

func TestWg(t *testing.T) {

	var wg sync.WaitGroup

	k := 0
	result := "res : "

	for i := 0; i < 20; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			k = k + 1
			result = result + "-"
			fmt.Printf("k : %d , i: %d, res : %s\n", k, i, result)
		}()
	}

	wg.Wait()
	fmt.Printf("k : %d , res : %s", k, result)

}

func TestGetResult(t *testing.T) {
	p := FofaPlugin{}
	res, _, _ := p.GetResult(1)
	log.Info(res)
}
