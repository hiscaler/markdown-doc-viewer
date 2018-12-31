package main

import (
	"config"
	"io/ioutil"
	"fmt"
	"path"
	"strings"
	"doc"
	"os"
	"encoding/json"
)

var cfg *config.Config

func init() {
	cfg = config.NewConfig()
}

func main() {
	docs := make([]doc.Doc, 0)
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
			docs = append(docs, doc)
		}
	}

	for _, doc := range docs {
		fmt.Println("Doc name is ", doc.Title)
		fmt.Println(fmt.Sprintf("%#v", doc))
	}
}
