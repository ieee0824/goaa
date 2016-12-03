package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

var list []string

func init() {
	infos, err := ioutil.ReadDir("static")
	if err != nil {
		os.Exit(1)
	}
	for _, info := range infos {
		list = append(list, info.Name())
	}
}
func root(w http.ResponseWriter, r *http.Request) {
}

func filelist(w http.ResponseWriter, r *http.Request) {
	bin, err := json.Marshal(list)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bin)
}

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/list", filelist)
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	http.ListenAndServe(":8080", nil)
}
