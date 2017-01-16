package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"os"
	"net/http/httputil"
	"net/url"
	//"reflect"
	//"regexp"
	"text/template"
	"github.com/codegangsta/negroni"
	"github.com/xyproto/mooseware"
	"github.com/PuerkitoBio/goquery"
)

var (
	//帖子路径正则表达式
	//threadItemExp = regexp.MustCompile(`"thread/[0123456789]+"`)
)

// This will be the index.html
var homeTemplate *template.Template

// This will store all the templates
var templates *template.Template

var arr []string

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
	//mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		arr := make([]string, 1)

		link := "http://tieba.baidu.com/f?kw=%D3%A2%D3%EF"
		doc, err := goquery.NewDocument(link)
		if err != nil {
			//logger.Error("解析页面失败：%s, %s", link, err.Error())
			log.Fatal(err)
		}

		doc.Find("li.j_thread_list.clearfix").Each(func(i int, s *goquery.Selection) {
			//title := s.Find(".threadlist_title").Text()
			//log.Println("第", i + 1, "个帖子的标题：", title)
			// 返回的是 *html.Node
			topicNode, err := s.Find(".threadlist_title").Html()
			if err != nil {
				log.Fatal(err)
			}
			log.Println(topicNode)
			arr = append(arr, topicNode)
		})

		data := struct {
			Title string
			TopicNodes []string
		}{
			Title: "golang html template demo",
			TopicNodes: arr,
		}
		//err = t.Execute(os.Stdout, data)
		err = homeTemplate.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
	})

	mux.HandleFunc("/p/", ProxyFunc)

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

func ProxyFunc(w http.ResponseWriter, r *http.Request) {
	// change the request host to match the target
	r.Host = "tieba.baidu.com"
	u, _ := url.Parse("http://tieba.baidu.com/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	// You can optionally capture/wrap the transport if that's necessary (for
	// instance, if the transport has been replaced by middleware). Example:
	// proxy.Transport = &myTransport{proxy.Transport}
	//proxy.Transport = &myTransport{}

	proxy.ServeHTTP(w, r)
}