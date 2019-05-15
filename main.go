package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/amendgit/kit"
)

var (
	sourceDir    = kit.SourceDir()
	cardsDir     = path.Join(sourceDir, "cards")
	metadataPath = path.Join(sourceDir, "metadata.json")
)

func main() {
	if len(os.Args) <= 1 {
		review()
		return
	}
	opt, args := os.Args[1], os.Args[2:]
	switch opt {
	case "help":
		help(args)
	case "u", "update":
		update(args)
	case "sync":
		sync()
	case "build":
		build()
	}
}

// help 显示帮助
func help(args []string) {
	fmt.Print("help message of card")
}

// update 编辑或者更新一张卡片。
func update(args []string) {
	cardID := args[0]
	cardPath := path.Join(cardsDir, cardID+".md")
	if !kit.IsPathExist(cardPath) {
		bs := GenerateEmptyCardContent(cardID)
		ioutil.WriteFile(cardPath, bs, 0666)
	}
	exec.Command("code", cardPath).Run()
}

// review 选取下一张需要复习的卡片，并进行复习。
func review() {
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
