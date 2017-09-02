package main

import (
	"log"
	"os"

	"github.com/ivanturianytsia/domainfinder"
	"github.com/ivanturianytsia/domainfinder/parallel"
	"github.com/ivanturianytsia/domainfinder/thesaurus"
)

type domain struct {
	name      string `json:"name,omitempty"`
	available bool   `json:"available,omitempty"`
	err       error  `json:"err,omitempty"`
	from      string `json:"from,omitempty"`
}

func findDomains(word string, output chan<- domain, done chan<- struct{}) {

	// synonyms
	thesaurus := thesaurus.NewBigHugeThesaurus(os.Getenv("BHT_APIKEY"))
	// sprinkle
	sprinkler := domainfinder.NewSprinkler()
	// coolify
	coolifier := domainfinder.NewCoolifier()
	// domainify
	domainifier := domainfinder.NewDomainifier()
	// available
	available := domainfinder.NewAvailable()

	// synonyms
	syns, err := thesaurus.Synonyms(word)
	if err != nil {
		log.Fatalf("Failed when looking for synonyms for %s\n%s\n", word, err.Error())
		return
	}
	syns = append(syns, word)

	p := parallel.NewStringParallelizer(syns)

	p.Do(func(v interface{}, i int) interface{} {
		s := sprinkler.Sprinkle(v.(string))
		c := coolifier.Coolify(s)
		d := domainifier.Domainify(c)
		a, err := available.Available(d)
		output <- domain{name: d, available: a, err: err, from: word}
		return a
	})
	done <- struct{}{}
}
