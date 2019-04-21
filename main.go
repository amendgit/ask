package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/amendgit/base"
)

var soruceDir = base.SourceDir()

func main() {
	if len(os.Args) == 0 {
		return
	}
	opt, args := os.Args[1], os.Args[2:]
	switch opt {
	case "help":
		showHelp(args)
	case "+", "edit":
		editCard(args)
	case "n", "next":
		nextCard(args)
	case "test":
		testCard(args)
	}
}

func showHelp(args []string) {
	fmt.Print("help message of card")
}

func editCard(args []string) {
	cardName := args[0]
	wd, _ := os.Getwd()
	cardPath := fmt.Sprintf("%v/cards/%v.md", wd, cardName)
	if base.IsPathExist(cardPath) {
		exec.Command("subl", cardPath).Run()
		return
	}
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
	data := map[string]string{"id": cardName}
	buf := bytes.NewBuffer(nil)
	tmpl.Execute(buf, data)
	ioutil.WriteFile(cardPath, buf.Bytes(), 0666)
	exec.Command("subl", cardPath).Run()
}

func nextCard(args []string) {
	syncCards()
}

func syncCards() {
	fileInfos, _ := ioutil.ReadDir("./cards")
	cards := make([]map[string]string)
	for _, fileInfo := range fileInfos {
		card := make(map[string]string)
		modTime := fileInfo.ModTime()
		fmt.Printf("name: %v mod: %v\n", fileInfo.Name(), modTime)
	}
	// file, _ := os.Create("./cards.yaml")
}

func testCard(args []string) {
	cardName := args[0]
	wd, _ := os.Getwd()
	cardPath := fmt.Sprintf("%v/cards/%v.md", wd, cardName)
	bs, err := ioutil.ReadFile(cardPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	seps := []string{"<!--front-->", "<!--back-->"}
	lines := base.Lines(string(bs))
	l, h, i := 0, 0, 0
	for h < len(lines) {
		for h < len(lines) && (i == len(seps) || !strings.Contains(lines[h], seps[i])) {
			h++
		}
		component := strings.Join(lines[l:h], "\n")
		fmt.Println(component)
		l, h, i = h+1, h+1, i+1
	}
}
