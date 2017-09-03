package domainfinder

import (
	"math/rand"
	"strings"
	"unicode"
)

const allowedChars = "abcdefghijklmnopqrstuvwxyz0123456789_-"

type Domainifier interface {
	Domainify(word string) string
}

type domainifier struct {
	tlds []string
}

func NewDomainifier() Domainifier {
	return domainifier{
		tlds: []string{"com", "net"},
	}
}

func (s domainifier) randTlds() string {
	return s.tlds[rand.Intn(len(s.tlds))]
}

func (s domainifier) Domainify(word string) string {
	word = strings.ToLower(word)
	var newWord []rune
	for _, r := range word {
		if unicode.IsSpace(r) {
			r = '-'
		}
		if !strings.ContainsRune(allowedChars, r) {
			continue
		}
		newWord = append(newWord, r)
	}
	return string(newWord) + "." + s.randTlds()
}
