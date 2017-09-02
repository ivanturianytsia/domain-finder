package domainfinder

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type Sprinkler interface {
	Sprinkle(word string) string
}

type sprinkler struct {
	transforms []string
}

func NewSprinkler() Sprinkler {
	return sprinkler{
		transforms: []string{
			"%s",
			"%sapp",
			"%ssite",
			"%stime",
			"%shq",
			"%scenter",
			"%sgalaxy",
			"%sland",
			"get%s",
			"go%s",
			"lets%s",
			"yesto%s",
			"super%s",
			"amazing%s",
			"start%s",
		},
	}
}

func (s sprinkler) Sprinkle(word string) string {
	t := s.transforms[rand.Intn(len(s.transforms))]
	return fmt.Sprintf(t, word)
}
