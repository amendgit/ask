package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/amendgit/base"
)

var (
	soruceDir = base.SourceDir()
	cardsPath = path.Join(soruceDir, "cards")
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
		syncCards()
	case "test":
		testCard(args)
	}
}

func showHelp(args []string) {
	fmt.Print("help message of card")
}

func editCard(args []string) {
	cardName := args[0]
	cardPath := path.Join(cardsPath, cardName)
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

type Card struct {
	Name       string     `json:"name"`
	ReviewTime *time.Time `json:"reviewTime,omitempty"`
	Level      int        `json:"level"`
}

func nextCard(args []string) {
	reviewCard := findNextCardToRview()
	fmt.Printf("reviewCard: %v\n", reviewCard.Name)
	cardPath := path.Join(cardsPath, reviewCard.Name)
	bs, _ := ioutil.ReadFile(cardPath)
	components := componentsFromCardContent(string(bs))
	fmt.Println(components[0])
	fmt.Println(components[1])
	var input string
	fmt.Scanf("%s", &input)
	fmt.Println(components[2])
	fmt.Printf("1.记得      2.不记得\n")
	var option int
	fmt.Scanf("%d", &option)
	bs, _ = ioutil.ReadFile(path.Join(soruceDir, "cards.json"))
	if option == 1 {
		reviewCard.Level = reviewCard.Level + 1
		duration := time.Duration(math.Exp(float64(reviewCard.Level))*24) * time.Hour
		reviewTime := time.Now().Local().Add(duration)
		reviewCard.ReviewTime = &reviewTime
	} else if option == 2 {
		duration := time.Duration(math.Exp(float64(reviewCard.Level))*24) * time.Hour
		reviewTime := time.Now().Local().Add(duration)
		reviewCard.ReviewTime = &reviewTime
	} else {
		return
	}
	cards := []Card{}
	json.Unmarshal(bs, &cards)
	for i, card := range cards {
		if card.Name == reviewCard.Name {
			cards[i] = *reviewCard
			break
		}
	}
	bs, _ = json.MarshalIndent(cards, "", "    ")
	ioutil.WriteFile(path.Join(soruceDir, "cards.json"), bs, 0666)
}

func findNextCardToRview() *Card {
	bs, _ := ioutil.ReadFile(path.Join(soruceDir, "cards.json"))
	cards := []Card{}
	json.Unmarshal(bs, &cards)
	var expiredCard *Card
	now := time.Now()
	for _, card := range cards {
		if card.ReviewTime == nil {
			continue
		}
		if card.ReviewTime.Before(now) {
			expiredCard = &card
			break
		}
	}
	if expiredCard != nil {
		return expiredCard
	}
	var newCards []Card
	for _, card := range cards {
		if card.Level == 0 {
			newCards = append(newCards, card)
		}
	}
	if len(newCards) == 0 {
		fmt.Println("暂时没有可以复习的卡片")
		return nil
	}
	fmt.Println("随机抽一张卡片")
	randomCard := newCards[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(newCards))]
	return &randomCard
}

func syncCards() {
	finfos, _ := ioutil.ReadDir(cardsPath)
	cards := []interface{}{}
	for _, finf := range finfos {
		card := map[string]string{}
		card["name"] = finf.Name()
		card["level"] = "0"
		cards = append(cards, card)
	}
	cardsJson, _ := os.Create(path.Join(soruceDir, "cards.json"))
	defer cardsJson.Close()
	bs, _ := json.MarshalIndent(cards, "", "    ")
	println(string(bs))
	cardsJson.Write(bs)
}

func testCard(args []string) {
	cardName := args[0]
	cardPath := path.Join(cardsPath, cardName)
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

func componentsFromCardContent(content string) []string {
	seps := []string{"<!--front-->", "<!--back-->"}
	lines := base.Lines(string(content))
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
