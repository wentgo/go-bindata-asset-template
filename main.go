package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	r := chi.NewRouter()

	//r.Use(middleware.Logger)
	r.Use(logging)

	r.Get("/", index)
	r.Get("/user/login", userLogin)
	r.Get("/error/404", error404)
	r.NotFound(error404)

	//http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	FileServer(r, "/assets", http.Dir("assets"))

	fmt.Println("Listening on localhost:8080 ...")
	http.ListenAndServe(":8080", r)
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(time.Now().Format("2006-01-02 03:04:05"), r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["Title"] = "Index"
	tpl.ExecuteTemplate(w, "home-index.gohtml", data)
}

func userLogin(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["Title"] = "User Login"
	data["Username"] = "Sam Wang"
	tpl.ExecuteTemplate(w, "user-login.gohtml", data)
}

func error404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	tpl.ExecuteTemplate(w, "error-404.gohtml", nil)
}
