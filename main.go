package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var tmpl *template.Template

func main() {
	tmpl = parseTemplates()
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", index)
	http.HandleFunc("/url", getURLResult)
	http.HandleFunc("/google", googleSearch)
	http.HandleFunc("/search", googleSearch)
	http.HandleFunc("/download", downloadPage)
	http.HandleFunc("/doDownload", download)
	http.ListenAndServe(":10203", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	datama := map[string]interface{}{"head": "go portal",
		"foot": "Develoed by Najy", "text": "Go Portal"}
	err := tmpl.ExecuteTemplate(w, "index.html", datama)
	if err != nil {
		println(err.Error())
	}
}

func downloadPage(w http.ResponseWriter, r *http.Request) {
	datama := map[string]interface{}{"head": "go portal",
		"foot": "Develoed by Najy", "text": "Go Portal"}
	err := tmpl.ExecuteTemplate(w, "download.html", datama)
	if err != nil {
		println(err.Error())
	}
}

func getURLResult(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("q")
	data, err := callURL(url)
	var res string
	if err != nil {
		res = err.Error()
	} else {
		res = string(data)
	}

	datama := map[string]interface{}{"head": "go portal",
		"foot": "Develoed by Najy", "text": res}
	tmpl.ExecuteTemplate(w, "index.html", datama)
}

func googleSearch(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("q")
	aurl := "https://google.com/search?q=" + url.PathEscape(text)
	var res string

	data, err := callURL(aurl)

	if err != nil {
		res = err.Error()
	} else {
		res = string(data)
	}

	datama := map[string]interface{}{"head": "google search",
		"foot": "Develoed by Najy", "text": res}
	tmpl.ExecuteTemplate(w, "google.html", datama)

}

func download(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("link")
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	downloadFile(ctx, w, url)
}

func parseTemplates() *template.Template {
	templ := template.New("")
	err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			_, err = templ.ParseFiles(path)
			if err != nil {
				log.Println(err)
			}
		}

		return err
	})

	if err != nil {
		panic(err)
	}
	return templ
}
