package main

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/amendgit/kit"
	yaml "gopkg.in/yaml.v3"
)

// CardDAO 封装了从数据库和文件系统读取Card的逻辑细节。
type CardDAO struct {
}

// NewCardDAO 获取CardDAO的实例
func NewCardDAO() *CardDAO {
	return new(CardDAO)
}

// Save 创建一个新的Card记录到数据库中
func (cardDAO *CardDAO) Save(card *Card) {
	// todo
}

// Get 从数据库中获取一个Card记录
func (cardDAO *CardDAO) Get(id string) *Card {
	db := GetAskDB()
	rows, err := db.Query("select id, title, question, answer, review_time from cards")
	if err != nil {
		log.Fatal(err)
	}
	if rows.Next() {
		var card Card
		rows.Scan(&card.ID, &card.Title, &card.Question, &card.Answer, &card.ReviewTime)
		return &card
	} else {
		return nil
	}
}

// GetAllCards 获取当前所有的卡片
func (cardDAO *CardDAO) GetAllCards() []Card {
	db := GetAskDB()
	rows, err := db.Query("select id, title, question, answer, review_time from cards")
	if err != nil {
		log.Fatal(err)
	}
	var cards []Card
	for rows.Next() {
		var card Card
		rows.Scan(&card.ID, &card.Title, &card.Question, &card.Answer, &card.ReviewTime)
		cards = append(cards, card)
	}
	return cards
}

// Add 添加一个卡片的数据到数据库中
func (cardDAO *CardDAO) Add(card *Card) {
	db := GetAskDB()
	tx, _ := db.Begin()
	insertStmt, err := tx.Prepare("insert or ignore into cards(id, title, question, answer, hash) values(?, ?, ?, ?, ?)")
	defer tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	// 如果记录已经存在，但是hash不一样，说明文件的内容发生了变化，需要刷新该条记录。
	updateStmt, err := tx.Prepare("update cards set id=?, title=?, question=?, answer=?, hash=? where id==? and hash!=?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = insertStmt.Exec(card.ID, card.Title, card.Question, card.Answer, card.Hash)
	if err != nil {
		log.Println(err)
	}
	_, err = updateStmt.Exec(card.ID, card.Title, card.Question, card.Answer, card.Hash, card.ID, card.Hash)
	if err != nil {
		log.Println(err)
	}
}

// Update 更新一个卡片的数据
func (cardDAO *CardDAO) Update(card *Card) {
	db := GetAskDB()
	tx, _ := db.Begin()
	updateStmt, err := tx.Prepare("update cards set title=?, question=?, answer=?, review_time=?, level=? where id==?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = updateStmt.Exec(card.Title, card.Question, card.Answer, card.ReviewTime, card.Level, card.ID)
	if err != nil {
		log.Fatal(err)
	}
}

// Delete 从数据库中删除一条记录
func (cardDAO *CardDAO) Delete(id string) {
	// todo
}

// PickOneCard 优先找是否有过期的卡片，没有的话再找一张新的卡片。
func (cardDAO *CardDAO) PickOneCard() *Card {
	card := cardDAO.PickOneOutdateCard()
	if card != nil {
		return card
	}
	return cardDAO.PickOneNewCard()
}

// PickOneNewCard 从新的卡片中，随机选取一张卡片。
func (cardDAO *CardDAO) PickOneNewCard() *Card {
	db := GetAskDB()
	rows, err := db.Query("select id, title, question, answer, level from cards where level == 0")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var cards []Card
	for rows.Next() {
		var card Card
		err := rows.Scan(&card.ID, &card.Title, &card.Question, &card.Answer, &card.Level)
		if err != nil {
			log.Fatal(err)
		}
		cards = append(cards, card)
	}
	if len(cards) == 0 {
		return nil
	}
	rand.Seed(time.Now().Unix())
	return &cards[rand.Intn(len(cards))]
}

// PickOneOutdateCard 从过期的卡片中，随机选取一张卡片。
func (cardDAO *CardDAO) PickOneOutdateCard() *Card {
	db := GetAskDB()
	rows, err := db.Query("select id, title, question, answer, review_time from cards where review_time < date('now') and level > 0")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var cards []Card
	for rows.Next() {
		var card Card
		err := rows.Scan(&card.ID, &card.Title, &card.Question, &card.Answer, &card.ReviewTime, &card.Level)
		if err != nil {
			log.Fatal(err)
		}
		cards = append(cards, card)
	}
	if len(cards) == 0 {
		return nil
	}
	rand.Seed(time.Now().Unix())
	return &cards[rand.Intn(len(cards))]
}

// ReadFile 从文件中读取文件
func (cardDAO *CardDAO) ReadFile(path string) *Card {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	return cardDAO.ParseString(string(bs))
}

// ParseString 从字符串中解析卡片数据
func (cardDAO *CardDAO) ParseString(s string) *Card {
	seps := []string{"<!--front-->", "<!--back-->"}
	lines := kit.Lines(string(s))
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
	var card Card
	yaml.Unmarshal([]byte(components[0]), &card)
	card.Question = components[1]
	card.Answer = components[2]
	card.ReviewTime = time.Now()
	hash := md5.Sum([]byte(s))
	card.Hash = hex.EncodeToString(hash[:])
	// 当标题为空时，用问题填充。
	if card.Title == "" {
		card.Title = kit.Lines(card.Question)[0]
	}
	return &card
}
