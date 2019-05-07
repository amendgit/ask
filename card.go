package main

import (
	"bytes"
	"strings"
	"text/template"
	"time"
)

// Metadata 表示一个卡片的元数据。
type Metadata struct {
	ID    string   `yaml:"id"`
	Title string   `yaml:"title"`
	Tags  []string `yaml:"tags"`
}

func (metadata *Metadata) String() string {
	s := ""
	s = s + "tags: " + "[" + strings.Join(metadata.Tags, ",") + "]\n"
	s = s + "title: " + metadata.Title + "\n"
	return s
}

// Card 表示一个卡片的实体类。
type Card struct {
	Metadata
	Question   string
	Answer     string
	ReviewTime time.Time
	Hash       string
}

// GenerateEmptyCardContent 生成一张空的卡片的内容
func GenerateEmptyCardContent(id string) []byte {
	tmpl, _ := template.New("card").Parse(
		`---
id: {{.id}}
title: null
tags:
    - null
---

<!--front-->
todo

<!--back-->
todo
`)
	data := map[string]string{"id": id}
	buf := bytes.NewBuffer(nil)
	tmpl.Execute(buf, data)
	return buf.Bytes()
}
