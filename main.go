package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/amendgit/kit"
)

var (
	sourceDir    = kit.SourceDir()
	cardsDir     = path.Join(sourceDir, "cards")
	metadataPath = path.Join(sourceDir, "metadata.json")
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
		sync()
	case "build":
		build()
	default:
		ReviewCard()
	}
}

func showHelp(args []string) {
	fmt.Print("help message of card")
}

func editCard(args []string) {
	cardID := args[0]
	cardPath := path.Join(cardsDir, cardID+".md")
	if !kit.IsPathExist(cardPath) {
		bs := GenerateEmptyCardContent(cardID)
		ioutil.WriteFile(cardPath, bs, 0666)
	}
	exec.Command("code", cardPath).Run()
}

type CardMetadata struct {
	Name       string     `json:"name"`
	ReviewTime *time.Time `json:"reviewTime,omitempty"`
	Level      int        `json:"level"`
}

// ReviewCard 选取下一张需要复习的卡片，并进行复习。
func ReviewCard() {
	cardDAO := NewCardDAO()
	card := cardDAO.PickOneCard()
	fmt.Printf("准备复习卡片: %v\n\n", card.ID)
	fmt.Printf("Question:\n %s\n\n", card.Question)

	var anyKey string
	fmt.Scanf("%s", &anyKey)

	fmt.Printf("Answer:\n %s\n\n", card.Answer)

	var option int
	for option <= 0 || option > 2 {
		fmt.Printf("1.记得            2.不记得\n")
		fmt.Scanf("%d", &option)
	}

	if option == 1 {
		card.Level = card.Level + 1
	}
}

func nextCard(args []string) {
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

func sync() {
	cardDAO := NewCardDAO()
	cardFileInfos, _ := ioutil.ReadDir(cardsDir)
	for _, cardFileInfo := range cardFileInfos {
		card := cardDAO.ReadFile("./cards/" + cardFileInfo.Name())
		cardDAO.Update(card)
	}
}

func componentsFromString(content string) []string {
	seps := []string{"<!--front-->", "<!--back-->"}
	lines := kit.Lines(string(content))
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
	exec.Command("go", "build", "github.com/amendgit/ask").Run()
}
