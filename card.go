package main

import (
	"bytes"
	"text/template"
	"time"
)

// Card 表示一个卡片的实体类。
type Card struct {
	ID         string   `yaml:"id"`
	Title      string   `yaml:"title"`
	Tags       []string `yaml:"tags"`
	Question   string
	Answer     string
	ReviewTime time.Time
	Hash       string
	Level      int
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
