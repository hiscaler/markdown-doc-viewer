package main

import (
	"config"
	"io/ioutil"
	"path"
	"strings"
	"doc"
	"os"
	"encoding/json"
	"net/http"
	"io"
	"log"
	"html/template"
)

var cfg *config.Config

func init() {
	cfg = config.NewConfig()
}

func getDocs() map[string]doc.Doc {
	docs := make(map[string]doc.Doc, 0)
	dirs, _ := ioutil.ReadDir(cfg.DocumentDir)
	for _, dir := range dirs {
		if dir.IsDir() {
			section := doc.Section{}
			doc := doc.Doc{}
			file, err := os.Open(path.Join(cfg.DocumentDir, dir.Name(), "conf.json"))
			if err == nil {
				jsonByte, err := ioutil.ReadAll(file)
				if err == nil {
					json.Unmarshal(jsonByte, &doc)
				}
				file.Close()
			}
			doc.ModifyDatetime = dir.ModTime()
			if len(doc.Title) == 0 {
				doc.Title = dir.Name()
			}
			if len(doc.Author) == 0 {
				doc.Author = "Unknown"
			}

			files, _ := ioutil.ReadDir(path.Join(cfg.DocumentDir, dir.Name()))
			for _, file := range files {
				if file.Name() != "conf.json" {
					filename := strings.TrimSuffix(path.Base(file.Name()), path.Ext(file.Name()))
					section.Title = filename
					f, _ := os.Open(path.Join(cfg.DocumentDir, dir.Name(), file.Name()))
					c, _ := ioutil.ReadAll(f)
					section.Content = string(c)
					doc.Sections = append(doc.Sections, section)
				}
			}
			docs[dir.Name()] = doc
		}
	}

	return docs
}

func RenderServer(w http.ResponseWriter, req *http.Request) {
	docs := getDocs()
	q := req.URL.Query()
	name := q.Get("name")
	if doc, ok := docs[name]; ok {
		//section := q.Get("section")
		t, err := template.ParseFiles("./view/index.html")
		if err != nil {
			log.Fatal(err)
		}
		err = t.Execute(w, doc)
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(w, "")
	} else {
		io.WriteString(w, name+" is not exist.")
	}

}

func main() {
	http.HandleFunc("/", RenderServer)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
