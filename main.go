package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"os"
	//"regexp"
	"text/template"
	"github.com/codegangsta/negroni"
	"github.com/xyproto/mooseware"
	//"github.com/PuerkitoBio/goquery"
)

var (
	//帖子路径正则表达式
	//threadItemExp = regexp.MustCompile(`"thread/[0123456789]+"`)
)

// This will be the index.html
var homeTemplate *template.Template

// This will store all the templates
var templates *template.Template

func main() {
	templates, err := template.ParseGlob("template/*.html")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// find the template with the name index.html
	homeTemplate = templates.Lookup("index.html")

	mux := http.NewServeMux()

	// 启动静态文件服务
	//mux.HandleFunc("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		//link := "http://tieba.baidu.com/f?kw=%D3%A2%D3%EF"
		//content, statusCode := httpGet(link)
		//doc, err := newDocumentFromURL(link)
//		doc, err := goquery.NewDocument(link)
//		if err != nil {
//			//logger.Error("解析页面失败：%s, %s", link, err.Error())
//			log.Fatal(err)
//		}
		//s := doc.Find("li.j_thread_list.clearfix").Text()
		//fmt.Fprint(w, s)
		//log.Fatal(s)

		data := struct {
			Title string
		}{
			Title: "golang html template demo",
		}
		//err = t.Execute(os.Stdout, data)
		err = homeTemplate.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
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