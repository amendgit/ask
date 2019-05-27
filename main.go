package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/amendgit/kit"
)

var (
	gSourceDir = kit.SourceDir()
	gCardsDir  = path.Join(gSourceDir, "cards")
)

func main() {
	db := GetAskDB()
	defer db.Close()
	if len(os.Args) <= 1 {
		review("")
		return
	}
	opt, args := os.Args[1], os.Args[2:]
	switch opt {
	case "help":
		help(args)
	case "w", "write":
		update(args)
	case "sync":
		sync()
	case "build":
		build()
	case "list":
		list()
	default:
		name := args[0]
		review(name)
	}
}

// help 显示帮助
func help(args []string) {
	fmt.Print("help message of card")
}

// update 编辑或者更新一张卡片。
func update(args []string) {
	cardID := args[0]
	cardPath := path.Join(gCardsDir, cardID+".md")
	if !kit.IsPathExist(cardPath) {
		bs := GenerateAnEmptyCard(cardID)
		ioutil.WriteFile(cardPath, bs, 0666)
	}
	exec.Command("code", cardPath).Run()
}

// review 选取下一张需要复习的卡片，并进行复习。
func review(name string) {
	cardDAO := NewCardDAO()
	var card *Card
	if len(name) >= 0 {
		card = cardDAO.Get(name)
	} else {
		card = cardDAO.PickOneCard()
	}
	fmt.Printf("准备复习卡片: %v\n\n", card.ID)
	fmt.Printf("Question:\n %s\n\n", card.Question)
	var anyKey string
	fmt.Scanf("%s", &anyKey)
	fmt.Printf("Answer:\n %s\n\n", card.Answer)
	var opt int
	for opt <= 0 || opt > 2 {
		fmt.Printf("是否记得（1.记得 2.不记得）: ")
		fmt.Scanf("%d", &opt)
	}
	if opt == 1 {
		card.Level = card.Level + 1
	}
	card.ReviewTime = time.Now().Add(time.Duration(math.Exp(float64(card.Level))*24) * time.Hour)
	cardDAO.Update(card)
}

// sync 将cards目录下的markdown文件，同步到数据库中去。
func sync() {
	cardDAO := NewCardDAO()
	cardFileInfos, _ := ioutil.ReadDir(gCardsDir)
	for _, cardFileInfo := range cardFileInfos {
		card := cardDAO.ReadFile("./cards/" + cardFileInfo.Name())
		if card != nil && !card.Draft {
			fmt.Printf("+ %-30v \t %-30v\n", card.ID, card.Title)
			cardDAO.Add(card)
		} else if card.Draft {
			fmt.Printf("~ %-30v \t %-30v \t ./cards/%v\n", card.ID, card.Title, cardFileInfo.Name())
		} else {
			fmt.Printf("? ./cards/%v\n", cardFileInfo.Name())
		}
	}
}

// list 用来列举当前有哪些问题
func list() {
	cardDAO := NewCardDAO()
	cards := cardDAO.GetAllCards()
	for i := 0; i < len(cards); i++ {
		card := cards[i]
		fmt.Printf("%-30v %v \t %v\n", card.ID, card.ReviewTime, card.Title)
	}
}

// build 方法构建自身
func build() {
	exec.Command("go", "install", "github.com/amendgit/ask").Run()
}
