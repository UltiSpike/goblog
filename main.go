package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.URL.Path == "/" {
		_, err := fmt.Fprint(w, "Hello~ love me")
		if err != nil {
			return
		}
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
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

func main() {
	router := mux.NewRouter()
	//router := http.NewServeMux()
	router.HandleFunc("/", homeHandler)
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
	http.ListenAndServe(":3000", router)

	homeURL, _ := router.Get("home").URL()
	fmt.Println("homeURL: ", homeURL)
	articleURL, _ := router.Get("articles.show").URL("id", "1")
	fmt.Println("articleURL: ", articleURL)

}
