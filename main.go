package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/amendgit/X"
	_ "github.com/mattn/go-sqlite3"
)

var (
	soruceDir    = X.SourceDir()
	cardsDir     = path.Join(soruceDir, "cards")
	metadataPath = path.Join(soruceDir, "metadata.json")
)

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
	case "sync":
		syncMetadataIfNeeded()
	case "test":
		testCard(args)
	case "build":
		build()
	}
}

func showHelp(args []string) {
	fmt.Print("help message of card")
}

func editCard(args []string) {
	cardName := args[0]
	cardPath := path.Join(cardsDir, cardName+".md")
	if X.IsPathExist(cardPath) {
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

type CardMetadata struct {
	Name       string     `json:"name"`
	ReviewTime *time.Time `json:"reviewTime,omitempty"`
	Level      int        `json:"level"`
}

func nextCard(args []string) {
	syncMetadataIfNeeded()
	bs, _ := ioutil.ReadFile(metadataPath)
	cardMetadatas := []CardMetadata{}
	json.Unmarshal(bs, &cardMetadatas)
	index := func() int {
		now := time.Now()
		for i := 0; i < len(cardMetadatas); i++ {
			if cardMetadatas[i].ReviewTime != nil && cardMetadatas[i].ReviewTime.Before(now) {
				return i
			}
		}
		var indexes []int
		for i := 0; i < len(cardMetadatas); i++ {
			if cardMetadatas[i].Level == 0 {
				indexes = append(indexes, i)
			}
		}
		if len(indexes) == 0 {
			return -1
		}
		fmt.Println("随机抽一张新的卡片学习")
		rand.Seed(time.Now().Unix())
		return indexes[rand.Intn(len(indexes))]
	}()
	if index == -1 {
		fmt.Println("暂时没有可以复习的卡片")
		return
	}
	cardMetadata := &cardMetadatas[index]
	fmt.Printf("准备复习卡片: %v\n", cardMetadata.Name)
	cardPath := path.Join(cardsDir, cardMetadata.Name)
	bs, _ = ioutil.ReadFile(cardPath)
	components := componentsFromString(string(bs))
	fmt.Println(components[0])
	fmt.Println(components[1])
	var anyKey string
	fmt.Scanf("%s", &anyKey)
	fmt.Printf("%v\n\n", components[2])
	var option int
	for option <= 0 || option > 1 {
		fmt.Printf("1.记得      2.不记得\n")
		fmt.Scanf("%d", &option)
	}
	if option == 1 {
		cardMetadata.Level = cardMetadata.Level + 1
	}
	bs, _ = ioutil.ReadFile(metadataPath)
	duration := time.Duration(math.Exp(float64(cardMetadata.Level))*24) * time.Hour
	reviewTime := time.Now().Local().Add(duration)
	cardMetadata.ReviewTime = &reviewTime
	bs, _ = json.MarshalIndent(cardMetadatas, "", "    ")
	ioutil.WriteFile(metadataPath, bs, 0666)
}

func syncMetadataIfNeeded() {
	db, err := sql.Open("sqlite3", "./ask.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(`
	create table if not exists cards (
		id integer primary key autoincrement,
		title char(50) not null,
		question text not null,
		answer text,
		create_time datetime default current_timestamp,
		review_time datetime,
		level integer default 0
	)`)
	if err != nil {
		log.Fatal(err)
	}

	metaInfo, _ := os.Stat(metadataPath)
	cardInfos, _ := ioutil.ReadDir(cardsDir)
	cards := []CardMetadata{}
	needUpdate := false
	if !X.IsPathExist(metadataPath) {
		needUpdate = true
	}
	for _, cardInfo := range cardInfos {
		card := CardMetadata{}
		card.Name = cardInfo.Name()
		card.Level = 0
		cards = append(cards, card)
		if !needUpdate && cardInfo.ModTime().After(metaInfo.ModTime()) {
			needUpdate = true
		}
	}
	if !needUpdate {
		return
	}
	fmt.Println("正在更新cards.json")
	bs, _ := json.MarshalIndent(cards, "", "    ")
	ioutil.WriteFile(metadataPath, bs, 0666)
}

func testCard(args []string) {
	cardName := args[0]
	cardPath := path.Join(cardsDir, cardName)
	bs, err := ioutil.ReadFile(cardPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	seps := []string{"<!--front-->", "<!--back-->"}
	lines := X.Lines(string(bs))
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

func componentsFromString(content string) []string {
	seps := []string{"<!--front-->", "<!--back-->"}
	lines := X.Lines(string(content))
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
	return components
}

func build() {
	goPath := os.Getenv("GOPATH")
	srcPath := path.Join(goPath, "src")
	pkgPath, _ := filepath.Rel(srcPath, soruceDir)
	exec.Command("go", "build", pkgPath).Run()
}
