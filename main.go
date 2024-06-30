package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
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

func articleStoreHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "创建新文章")
	if err != nil {
		return
	}
}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {

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
	router.HandleFunc("/articles", articleStoreHandler).Methods(
		"POST").Name("articles.store")
	// gorilla/mux 限定类型的方式 [0-9]+
	router.HandleFunc("/articles/{id:[0-9]+}", articleShowHandler).Methods(
		"GET").Name("articles.show")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")

	router.Use(forceHtmlMiddleware)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
