package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
	"strings"
	"text/template"
	"time"

	"github.com/amendgit/X"
	"gopkg.in/yaml.v2"
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

// FromFile 从文件中读取卡片内容并初始化。
func (card *Card) FromFile(path string) error {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return card.FromString(string(bs))
}

// FromString 从s中解析卡片数据并初始化。
func (card *Card) FromString(s string) error {
	seps := []string{"<!--front-->", "<!--back-->"}
	lines := X.Lines(string(s))
	var components []string
	l, h, i := 0, 0, 0
	for h < len(lines) {
		for h < len(lines) && (i == len(seps) || !strings.Contains(lines[h], seps[i])) {
			h++
		}
		component := strings.Join(lines[l:h], "\n")
		components = append(components, component)
		l, h, i = h+1, h+1, i+1
	}
	yaml.Unmarshal([]byte(components[0]), &card.Metadata)
	card.Question = components[1]
	card.Answer = components[2]
	card.ReviewTime = time.Now()
	hash := md5.Sum([]byte(s))
	card.Hash = hex.EncodeToString(hash[:])
	log.Printf("card %v", card.ID)
	return nil
}

// ToFile 将卡片的数据写入到文件中去。
func (card *Card) ToFile(path string) error {
	// todo
	return nil
}

// ToString 将卡片数据表示为Markdown格式的文本。
func (card *Card) ToString() (string, error) {
	// todo
	return "", nil
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
