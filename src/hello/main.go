package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	//"regexp"
	"github.com/codegangsta/negroni"
	"github.com/xyproto/mooseware"
	"github.com/PuerkitoBio/goquery"
)

var (
	//帖子路径正则表达式
	//threadItemExp = regexp.MustCompile(`"thread/[0123456789]+"`)
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		link := "http://tieba.baidu.com/f?kw=%D3%A2%D3%EF"
		//content, statusCode := httpGet(link)
		//doc, err := newDocumentFromURL(link)
		doc, err := goquery.NewDocument(link)
		if err != nil {
			//logger.Error("解析页面失败：%s, %s", link, err.Error())
			//return nil
		}
		s := doc.Find("li.j_thread_list.clearfix").Text()
		fmt.Fprint(w, s)
	})

	n := negroni.Classic()

	// Moose status
	n.Use(moose.NewMiddleware())

	// Handler goes last
	n.UseHandler(mux)

	// Serve
	n.Run(":3000")
}

func httpGet(url string) (content string, statusCode int) {
	res, err1 := http.Get(url)
	if err1 != nil {
		statusCode = -100
		return
	}
	defer res.Body.Close()
	data, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		statusCode = -200
		return
	}
	statusCode = res.StatusCode
	content = string(data)
	return
}