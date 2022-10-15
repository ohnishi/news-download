package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"text/template"
	"time"

	"github.com/ohnishi/nahaha/backend/cmd"
	"github.com/pkg/errors"
)

func publishTrends(src, dest string, date time.Time) (err error) {
	srcPath := filepath.Join(src, date.Format("20060102")+".json")
	f, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}

	var content cmd.Content
	if err = json.Unmarshal(f, &content); err != nil {
		return err
	}

	if len(content.Items) == 0 {
		return errors.New("content size is zero")
	}

	if err = writeContent(dest, date, content); err != nil {
		return err
	}

	return nil
}

func writeContent(dest string, date time.Time, content cmd.Content) error {
	f, err := cmd.CreateOutFile(filepath.Join(dest, date.Format("2006/01/02")+".md"))
	if err != nil {
		return err
	}
	defer f.Close()

	funcMap := template.FuncMap{
		"rank": func(a int) int { return a + 1 },
	}
	t := template.Must(template.New("funcmap").Funcs(funcMap).Parse(tmplStr))

	// Execute(io.Writer(出力先), データ)
	if err := t.Execute(f, content); err != nil {
		log.Fatal(err)
	}

	if err := f.Sync(); err != nil {
		return errors.Wrap(err, "failed to sync file")
	}
	return nil
}

const tmplStr = `
---
title: "{{ .FormatDate }} に話題になった人"
date: {{ .Date }}
sidebar: "right"
---

{{ range $i, $item := .Items -}}
### {{ rank $i }}位 {{ $item.Word }} （{{ $item.Count }}記事）
{{ range $j, $article := $item.Articles -}}
- [{{ $article.Title }}]({{ $article.URL }})
{{ end }}
{{ end }}
`
