package main

import (
	"fmt"
	"net/http"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.URL.Path == "/" {
		_, err := fmt.Fprint(w, "Hello~ love me")
		if err != nil {
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "no pages")
	}
}

func otherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := fmt.Fprint(w, "此博客是用以记录编程笔记，请联系"+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
	if err != nil {
		return
	}
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", defaultHandler)
	router.HandleFunc("/about", otherHandler)
	http.ListenAndServe(":3000", router)

}
