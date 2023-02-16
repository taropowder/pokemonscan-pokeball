package utils

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func GetHeaderStr(header http.Header) (res string) {
	for key, value := range header {
		res = res + fmt.Sprintf("%s:%s\n", key, value)
	}
	return
}

func GetUrlInfo(url string) (res string, statusCode int, respLength int, err error) {
	//time.Sleep(time.Second * 4)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{Timeout: 10 * time.Second, Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	res = Md5(string(body))
	statusCode = resp.StatusCode
	respLength = len(string(body))
	return
}
