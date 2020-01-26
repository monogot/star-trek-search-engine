package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/davecgh/go-spew/spew"
)

type Dialog []struct {
	Episodename string `json:"episodename"`
	Act         string `json:"act"`
	Scenenumber string `json:"scenenumber"`
	Texttype    string `json:"texttype"`
	Who         string `json:"who"`
	Text        string `json:"text"`
	Speech      string `json:"speech"`
	Released    string `json:"released"`
	Episode     string `json:"episode"`
	Imdbrating  string `json:"imdbrating"`
	Imdbid      string `json:"imdbid"`
	Season      string `json:"season"`
}

var tpl = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	params := u.Query()
	searchKey := params.Get("q")

	dialog := &Dialog{}

	endpoint := fmt.Sprintf("https://tngapi-awicwils6q-ew.a.run.app/?q=%s", url.QueryEscape(searchKey))
	resp, err := http.Get(endpoint)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&dialog)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	spew.Dump(dialog)

	err = tpl.Execute(w, dialog)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":3000", mux)
}
