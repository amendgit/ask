package main

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/amendgit/X"
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
	// todo
	return nil
}

// Update 更新一个卡片的数据到数据库中
func (cardDAO *CardDAO) Update(card *Card) {
	db := GetAskDB()
	tx, _ := db.Begin()
	insertStmt, err := tx.Prepare("insert or ignore into cards(id, title, question, answer, hash) values(?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
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
	tx.Commit()
}

// Delete 从数据库中删除一条记录
func (cardDAO *CardDAO) Delete(id string) {

}

// ReadFromFile 从文件中读取文件
func (cardDAO *CardDAO) ReadFromFile(path string) *Card {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	return cardDAO.ParseFromString(string(bs))
}

// ParseFromString 从字符串中解析卡片数据
func (cardDAO *CardDAO) ParseFromString(s string) *Card {
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
	var card Card
	yaml.Unmarshal([]byte(components[0]), &card.Metadata)
	card.Question = components[1]
	card.Answer = components[2]
	card.ReviewTime = time.Now()
	hash := md5.Sum([]byte(s))
	card.Hash = hex.EncodeToString(hash[:])
	return &card
}
