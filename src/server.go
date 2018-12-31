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
			doc := doc.Doc{
				Name:           dir.Name(),
				ModifyDatetime: dir.ModTime(),
				Sections:       make(map[string]doc.Section),
			}
			file, err := os.Open(path.Join(cfg.DocumentDir, dir.Name(), "conf.json"))
			if err == nil {
				jsonByte, err := ioutil.ReadAll(file)
				if err == nil {
					json.Unmarshal(jsonByte, &doc)
				}
				file.Close()
			}
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
					doc.Sections[filename] = section
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
	if d, ok := docs[name]; ok {
		t, err := template.ParseFiles("./template/index.html")
		if err != nil {
			log.Fatal(err)
		}

		section := q.Get("section")
		content := ""
		if section, ok := d.Sections[section]; !ok {
			for _, section := range d.Sections {
				// Get first section content
				content = section.Content
				break
			}
		} else {
			content = section.Content
		}

		docItems := map[string]string{}
		for docName, t := range docs {
			docItems[docName] = t.Name
		}
		err = t.Execute(w, struct {
			Docs    map[string]string
			Doc     doc.Doc
			Content string
		}{
			Docs:    docItems,
			Doc:     d,
			Content: content,
		})
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(w, "")
	} else {
		io.WriteString(w, name+" is not exist.")
	}
}

func main() {
	http.Handle("/css/", http.FileServer(http.Dir("template")))
	http.HandleFunc("/", RenderServer)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
