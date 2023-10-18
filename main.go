package main

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"qst-ext-appsearcher-go/config"
	"qst-ext-appsearcher-go/pb/defs"
	"qst-ext-appsearcher-go/pb/extension"
	"qst-ext-appsearcher-go/trie"
	"strings"
)

type server struct {
	extension.UnimplementedExtInteractServer
	Trie *trie.Trie
}

var Apps = trie.NewTrie()

func (s *server) Search(_ context.Context, in *extension.Input) (*extension.SearchResult, error) {
	vs := Apps.StartsWith(in.Content)
	dl := extension.DisplayList{}
	for _, v := range vs {
		aid := v.(uint32)
		ai := config.Status.Attr.Apps[aid]
		dl.List = append(dl.List, &extension.DisplayItem{
			ObjId: aid,
			Name:  ai.Name,
			Hint:  nil,
		})
	}
	return &extension.SearchResult{Mresult: &extension.SearchResult_Ok{Ok: &extension.SearchResult_MOk{DisplayList: &dl}}}, nil
}
func (s *server) Submit(_ context.Context, hint *extension.SubmitHint) (*defs.MResult, error) {
	ai := config.Status.Attr.Apps[hint.ObjId]
	argv := append([]string{}, ai.Exec)
	attr := &os.ProcAttr{}
	attr.Dir = ai.Dir
	h, err := os.StartProcess(ai.Name, argv, attr)
	if err != nil {
		return &defs.MResult{Mresult: &defs.MResult_Status{Status: &defs.Status{Type: 1}}}, nil
	}
	config.Status.RunStat[hint.ObjId] = &config.App{Process: h}
	return &defs.MResult{Mresult: &defs.MResult_Ok{Ok: &defs.MResult_MOk{}}}, nil
}

func update() {
	cdir := "/usr/share/applications"
	if err := filepath.WalkDir(cdir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// 处理遍历过程中的错误
			return err
		}
		if path == "/usr/share/applications" {
			return nil
		}
		log.Println("prase:", path)
		if d.IsDir() {
			log.Println("skip:", path)
			return filepath.SkipDir
		}
		//read file content with path
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		var start = false
		for scanner.Scan() {
			var line = scanner.Text()
			if line == "[Desktop Entry]" {
				start = true
				break
			}
		}
		if start {
			var app = config.AppAttr{}
			for scanner.Scan() {
				var line = scanner.Text()
				//startwith
				if strings.HasPrefix(line, "[") {
					break
				} else if strings.HasPrefix(line, "Name=") {
					app.Name = line[5:]
				} else if strings.HasPrefix(line, "Exec=") {
					app.Exec = line[5:]
				} else if strings.HasPrefix(line, "Icon=") {
					app.Icon = line[5:]
				}
			}
			log.Printf("app: %+v\n", app)
			Apps.Insert(app.Name, app)
		}
		if err := scanner.Err(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		fmt.Println("Error:", err)
	}
}

func main() {
}
