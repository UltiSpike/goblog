package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		_, err := fmt.Fprint(w, "Hello~ love me")
		if err != nil {
			return
		}
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

func articleShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, i := fmt.Fprintf(w, "文章id"+id)
	if i != nil {
		return
	}
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "访问文章列表")
	if err != nil {
		return
	}
}

type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {

	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := make(map[string]string)

	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}

	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}

	// 检查是否有错误
	if len(errors) == 0 {
		fmt.Fprint(w, "验证通过!<br>")
		fmt.Fprintf(w, "title 的值为: %v <br>", title)
		fmt.Fprintf(w, "title 的长度为: %v <br>", utf8.RuneCountInString(title))
		fmt.Fprintf(w, "body 的值为: %v <br>", body)
		fmt.Fprintf(w, "body 的长度为: %v <br>", utf8.RuneCountInString(body))
	} else {

		storeURL, _ := router.Get("articles.store").URL()

		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
	}
}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {

	storeURL, _ := router.Get("articles.store").URL()
	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}
	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func forceHtmlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// 包装router
func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 除首页外，移除所有请求路径后面的斜杠
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimRight(r.URL.Path, "/")
		}

		// 传递请求
		next.ServeHTTP(w, r)
	})
}

var router = mux.NewRouter()

func main() {
	//router := http.NewServeMux()

	router.HandleFunc("/", homeHandler).Name("home")
	router.HandleFunc("/about", aboutHandler)
	router.HandleFunc("/articles/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.SplitN(r.URL.Path, "/", 3)[2]
		fmt.Fprintf(w, "文章id "+id)
	})
	router.HandleFunc("/articles", articlesIndexHandler).Methods(
		"GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods(
		"POST").Name("articles.store")
	// gorilla/mux 限定类型的方式 [0-9]+
	router.HandleFunc("/articles/{id:[0-9]+}", articleShowHandler).Methods(
		"GET").Name("articles.show")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")

	router.Use(forceHtmlMiddleware)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
