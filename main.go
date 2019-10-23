package main

//go:generate go-bindata -fs assets/js assets/css assets/img templates

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

const DEBUG = false

var tpl *template.Template

func init() {
	if DEBUG {
		// load templates from local file system
		tpl = template.Must(template.ParseGlob("templates/*"))
		return
	}

	// load templates from bindata
	var allTemplates string

	files, err := AssetDir("templates")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		content, err := Asset("templates/" + file)
		if err != nil {
			panic(err)
		}
		allTemplates += `{{define "` + file + `"}}` + string(content) + "{{end}}"
	}

	tpl = template.Must(template.New("root").Parse(allTemplates))
}

func main() {
	r := chi.NewRouter()

	r.Use(logging) // r.Use(middleware.Logger)

	r.Get("/", index)
	r.Get("/user/login", userLogin)
	r.Get("/error/404", error404)

	r.NotFound(error404)

	//http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	r.Get("/assets/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fs := http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))
		//fs := http.StripPrefix("/assets/", http.FileServer(AssetFile()))
		fs := http.FileServer(AssetFile())
		fs.ServeHTTP(w, r)
	}))

	fmt.Println("Listening on localhost:8080 ...")
	http.ListenAndServe(":8080", r)
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
