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
	"github.com/russross/blackfriday"
	"strconv"
	"hash/fnv"
)

var cfg *config.Config

func init() {
	cfg = config.NewConfig()
}

func getDocs() map[string]doc.Doc {
	docs := make(map[string]doc.Doc, 0)
	dirs := make([]string, 0)
	if ds, err := ioutil.ReadDir(cfg.DocumentDir); err == nil {
		for _, d := range ds {
			if d.IsDir() {
				dirs = append(dirs, path.Join(cfg.DocumentDir, d.Name()))
			}
		}
	}

	for _, dir := range cfg.DocumentDirs {
		dir = path.Clean(dir)
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			dirs = append(dirs, dir)
		}
	}

	h := fnv.New64()
	for _, dirPath := range dirs {
		if dir, err := os.Stat(dirPath); err == nil {
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
			files, _ := ioutil.ReadDir(dirPath)
			for _, file := range files {
				if file.Name() != "conf.json" {
					filename := strings.TrimSuffix(path.Base(file.Name()), path.Ext(file.Name()))
					section.Title = filename
					f, _ := os.Open(path.Join(dirPath, file.Name()))
					c, _ := ioutil.ReadAll(f)
					section.Content = c
					doc.Sections[filename] = section
				}
			}

			h.Reset()
			h.Write([]byte(dirPath))
			docs[strconv.FormatUint(h.Sum64(), 10)] = doc
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
		content := make([]byte, 0)
		if s, ok := d.Sections[section]; !ok {
			for k, v := range d.Sections {
				// Get first section content
				section = k
				content = v.Content
				break
			}
		} else {
			content = s.Content
		}

		if len(content) > 0 {
			content = blackfriday.MarkdownCommon(content)
		}

		docItems := map[string]string{}
		for docName, t := range docs {
			docItems[docName] = t.Name
		}
		err = t.Execute(w, struct {
			CurrentDocName     string
			CurrentSectionName string
			Docs               map[string]string
			Doc                doc.Doc
			Content            template.HTML
		}{
			CurrentDocName:     name,
			CurrentSectionName: section,
			Docs:               docItems,
			Doc:                d,
			Content:            template.HTML(string(content)),
		})
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(w, "")
	} else {
		t, err := template.ParseFiles("./template/list.html")
		if err != nil {
			log.Fatal(err)
		}
		docItems := make(map[string]string, 0)
		for k, v := range docs {
			docItems[k] = v.Name
		}
		err = t.Execute(w, struct {
			Docs map[string]string
		}{
			Docs: docItems,
		})
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(w, "")
	}
}

func main() {
	http.Handle("/css/", http.FileServer(http.Dir("template")))
	http.HandleFunc("/", RenderServer)
	err := http.ListenAndServe(":"+strconv.Itoa(cfg.Port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
