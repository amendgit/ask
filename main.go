package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

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

// sync 将cards目录下的markdown文件，同步到数据库中去。
func sync() {
	cardDAO := NewCardDAO()
	cardFileInfos, _ := ioutil.ReadDir(cardsDir)
	for _, cardFileInfo := range cardFileInfos {
		card := cardDAO.ReadFile("./cards/" + cardFileInfo.Name())
		cardDAO.Update(card)
	}
}

// build 方法构建自身
func build() {
	exec.Command("go", "build", "github.com/amendgit/ask").Run()
}
