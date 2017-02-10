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
		req.Header.Add("Cookie", `TIEBA_USERTYPE=35f770fc37174a4733210804; TIEBAUID=4744b33b40f25e783e2053ed; bdshare_firstime=1461033889439; rpln_guide=1; BDUSS=J5eFpVTjc3cTJjWmwwcy1SWTkwY3k2T3VYZ0w4QUU1UFpIYnRZUnhsTjVqQ0JZQVFBQUFBJCQAAAAAAAAAAAEAAAA5DsMNwau7qLXE7r-w7QAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAHn~-Fd5~~hXZl; __cfduid=d49e1554629f98e985e49b4211d0e4ed01480938459; Hm_lvt_287705c8d9e2073d13275b18dbd746dc=1480150516,1481251799; BAIDUID=790368EB0B68D1F4C99B5CE8D201AA79:FG=1; PSTM=1483019965; BIDUPSID=F76E100A4FDCF8E145D344549AEDE87D; batch_delete_mode=true; STOKEN=7cdeb94b92fc441e4ebbeeaeeb4ce65c00d5ffbc16dd6b45c7ca9b3cb7cf7986; BDSFRCVID=JfKsJeCCxG36zJ6iqCwFOwAo_P_-JdjDIBG33J; H_BDCLCKID_SF=JJkO_D_atKvjDbTnMITHh-F-5fIX5-RLf2Qha-OF5lOTJh0RjR6kj4FLLpJlWbbCaaOTM56F5ncUOqojDTbke4tX-NFetjFJJf5; wise_device=0; BDRCVFR[feWj1Vr5u3D]=I67x6TjHwwYf0; PSINO=5; H_PS_PSSID=21871_1422_21086_21803`)
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
			//log.Println(topicNode)
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
	http.Redirect(w, r, "http://tieba.baidu.com" + r.URL.String(), http.StatusFound)

	// change the request host to match the target
	//r.Host = "tieba.baidu.com"
	//u, _ := url.Parse("http://tieba.baidu.com/")
	//proxy := httputil.NewSingleHostReverseProxy(u)
	// You can optionally capture/wrap the transport if that's necessary (for
	// instance, if the transport has been replaced by middleware). Example:
	// proxy.Transport = &myTransport{proxy.Transport}
	//proxy.Transport = &myTransport{}

	//proxy.ServeHTTP(w, r)
}