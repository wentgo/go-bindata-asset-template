package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	http.HandleFunc("/", logging(index))
	http.HandleFunc("/user/login", logging(userLogin))
	http.HandleFunc("/error/404", logging(error404))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	fmt.Println("Listening on localhost:8080 ...")
	http.ListenAndServe(":8080", nil)
}

func logging(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(time.Now().Format("2006-01-02 03:04:05"), r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		//http.Redirect(w, r, "/", http.StatusSeeOther)
		//http.Redirect(w, r, "/error/404", 301)
		//http.NotFound(w, r)
		error404(w, r)
		return
	}
	data := make(map[string]interface{})
	data["Title"] = "Index"
	data["User"] = "Sam Wang"
	tpl.ExecuteTemplate(w, "home-index.gohtml", data)
}

func userLogin(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["Title"] = "User Login"
	tpl.ExecuteTemplate(w, "user-login.gohtml", data)
}

func error404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	tpl.ExecuteTemplate(w, "error-404.gohtml", nil)
}
