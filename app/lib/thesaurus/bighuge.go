package thesaurus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BigHuge struct {
	APIKey string
}

type synonyms struct {
	Noun *words `json:"noun"`
	Verb *words `json:"verb"`
}
type words struct {
	Syn []string `json:"syn"`
}

func NewBigHugeThesaurus(APIKey string) Thesaurus {
	return BigHuge{
		APIKey: APIKey,
	}
}

func (b BigHuge) Synonyms(term string) ([]string, error) {
	var syns []string

	response, err := http.Get("http://words.bighugelabs.com/api/2/" + b.APIKey + "/" + term + "/json")
	if err != nil {
		return syns, fmt.Errorf("BigHuge: Failed when looking for synonyms for %s:\n%s", term, err.Error())
	}

	var data synonyms

	defer response.Body.Close()
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return syns, err
	}
	if data.Noun != nil {
		syns = append(syns, data.Noun.Syn...)
	}
	if data.Verb != nil {
		syns = append(syns, data.Verb.Syn...)
	}
	return syns, nil
}
