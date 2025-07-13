package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))
	http.HandleFunc("/api/add", handleAdd)
	http.HandleFunc("/api/csv/", handleCSV)
	fmt.Println("Server running at http://localhost:8080/web/")
	http.ListenAndServe(":8080", nil)
}

func handleCSV(w http.ResponseWriter, r *http.Request) {
	// e.g. /api/csv/books.csv
	file := strings.TrimPrefix(r.URL.Path, "/api/csv/")
	f, err := os.Open("resources/" + file)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	defer f.Close()
	io.Copy(w, f)
}

type AddRequest struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	var req AddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(400)
		return
	}
	var path string
	var record []string
	switch req.Type {
	case "book":
		path = "resources/books.csv"
		record = []string{
			req.Data["title"], req.Data["isbn"], req.Data["authors"], req.Data["description"],
		}
	case "magazine":
		path = "resources/magazines.csv"
		record = []string{
			req.Data["title"], req.Data["isbn"], req.Data["authors"], req.Data["publishedAt"],
		}
	case "author":
		path = "resources/authors.csv"
		record = []string{
			req.Data["email"], req.Data["firstName"], req.Data["lastName"],
		}
	default:
		w.WriteHeader(400)
		return
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer f.Close()
	wtr := csv.NewWriter(f)
	wtr.Comma = ';'
	if err := wtr.Write(record); err != nil {
		w.WriteHeader(500)
		return
	}
	wtr.Flush()
	w.WriteHeader(200)
}
