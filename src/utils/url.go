package utils

import (
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
	"time"
)

func GetHeaderStr(header http.Header) (res string) {
	for key, value := range header {
		res = res + fmt.Sprintf("%s:%s\n", key, value)
	}
	return
}

func GetUrlInfo(url string) (pageHash string, statusCode int, title string, respLength int, err error) {
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
	// HTML 解析 body,并提取其中所有 css 以及 js 链接
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		b := make([]byte, 0)
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		respLength = len(b)
		return
	}
	// 获取 title

	title = doc.Find("title").Text()

	var cssLinks []string
	var jsLinks []string

	// 获取所有的 CSS 链接
	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		if rel, _ := s.Attr("rel"); rel == "stylesheet" {
			if href, ok := s.Attr("href"); ok {
				cssLinks = append(cssLinks, href)
			}
		}
	})

	// 获取所有的 JavaScript 链接
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		if src, ok := s.Attr("src"); ok {
			jsLinks = append(jsLinks, src)
		}
	})
	hashString := doc.Text()
	if len(cssLinks) != 0 || len(jsLinks) != 0 {
		hashString = strings.Join(cssLinks, "|") + strings.Join(jsLinks, "|")
	}
	pageHash = Md5(hashString)

	statusCode = resp.StatusCode
	respLength = len(hashString)
	return
}
