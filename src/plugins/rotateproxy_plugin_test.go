package plugins

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"testing"
)

func TestRotateProxy(t *testing.T) {
	r := RotateProxyPlugin{}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("获取当前工作目录失败:", err)
		return
	}

	fmt.Println("当前工作目录:", cwd)

	filePath := "data/config/rotate_proxy.json"
	configStr, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("读取文件失败:", err)
		return
	}
	//configStr := "{\n  \"whitelist_type\": [\n    \"jpg\",\n    \"gif\",\n    \"png\",\n    \"css\"\n  ],\n  \"request_intercept_rules\": {\n    \"uri\": {\n      \"login_page\": \"login\"\n    },\n    \"headers\": {\n      \"Shiro\": \"deleteMe\"\n    },\n    \"parameters\": {\n      \"password_filed\": \"passwd\"\n    }\n  },\n  \"response_intercept_rules\": {\n    \"uri\": {\n    },\n    \"headers\": {\n      \"Shiro\": \"deleteMe\"\n    },\n    \"body\": {\n      \"password_filed\": \"passwd\"\n    }\n  }\n}"
	//configStr := "{\n  \"resp_intercept_rules\": [\n    {\n      \"url\": [\"js\"],\n      \"data\": [\"{path:\"]\n    }\n  ]\n}"
	//configStr := "{\n  \"resp_intercept_rules\": [\n    {\n      \"url\": [\"js\"],\n      \"data\": [\"this.message,name\"]\n    }\n  ],\n  \"port\": 8980\n}"
	r.Register(nil, string(configStr))

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)

	<-sigc

	os.Exit(0)

	//r, _ := regexp.Compile("((cmd=)|(exec=)|(command=)|(execute=)|(ping=)|(query=)|(jump=)|(code=)|(reg=)|(do=)|(func=)|(arg=)|(option=)|(load=)|(process=)|(step=)|(read=)|(function=)|(feature=)|(exe=)|(module=)|(payload=)|(run=)|(daemon=)|(upload=)|(dir=)|(download=)|(log=)|(ip=)|(cli=))")
	//res := r.FindString("{\"csp-report\":{\"document-uri\":\"https://share.doppler.com/\",\"referrer\":\"\",\"violated-directive\":\"script-src\",\"effective-directive\":\"script-src\",\"original-policy\":\"upgrade-insecure-requests;default-src 'none';script-src 'self' 'unsafe-inline' https://static.cloudflareinsights.com 'nonce-Ux/86p5pOXqJnPYPtDRVHG7epom9iWZycYWjFO782NQ';style-src 'self' 'unsafe-inline';img-src 'self' data: https://doppler.com;connect-src 'self';font-src 'self' data:;form-action 'self';frame-ancestors 'none';base-uri 'self';report-uri https://doppler.report-uri.com/r/d/csp/enforce\",\"disposition\":\"enforce\",\"blocked-uri\":\"eval\",\"line-number\":3,\"column-number\":155,\"status-code\":200,\"script-sample\":\"\"}}")
	//if res != "" {
	//	fmt.Println(res) //Hello World!
	//} else {
	//	fmt.Println("null")
	//}

}
