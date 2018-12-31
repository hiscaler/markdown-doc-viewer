package main

import (
	"config"
	"io/ioutil"
	"fmt"
	"path"
	"strings"
)

var cfg *config.Config

func init() {
	cfg = config.NewConfig()
}

func main() {
	docs := make([]map[string]string, 0)
	dirs, _ := ioutil.ReadDir(cfg.DocumentDir)
	for _, dir := range dirs {
		if dir.IsDir() {
			mdDocs := make(map[string]string)
			files, _ := ioutil.ReadDir(path.Join(cfg.DocumentDir, dir.Name()))
			for _, file := range files {
				filename := strings.TrimSuffix(path.Base(file.Name()), path.Ext(file.Name()))
				mdDocs[filename] = path.Join(cfg.DocumentDir, dir.Name(), file.Name())
			}
			docs = append(docs, mdDocs)
		}
	}

	fmt.Println(fmt.Sprintf("%#v", docs))
}
