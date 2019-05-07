package main

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"strings"
	"time"

	"github.com/amendgit/X"
	"gopkg.in/yaml.v2"
)

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

// Card represent an data of card.
type Card struct {
	Metadata   Metadata
	ID         string
	Question   string
	Answer     string
	ReviewTime time.Time
	Hash       string
}

// FromFile init card from string.
func (card *Card) FromFile(path string) error {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return card.FromString(string(bs))
}

// FromString init card from string.
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
	return nil
}

// ToFile write card data to file as markdown format.
func (card *Card) ToFile(path string) error {
	// todo
	return nil
}

// ToString write card data to string as markdown format.
func (card *Card) ToString() (string, error) {
	// todo
	return "", nil
}
