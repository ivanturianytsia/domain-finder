package domainfinder

const (
	duplicateVowel bool = true
	removeVowel    bool = false
)

type Coolifier interface {
	Coolify(word string) string
}

type coolifier struct{}

func NewCoolifier() Coolifier {
	return coolifier{}
}

func (c coolifier) Coolify(word string) string {
	wordB := []byte(word)

	if randBool() {
		var vI int = -1
		for i, char := range wordB {
			switch char {
			case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
				if randBool() {
					vI = i
				}
			}
		}
		if vI >= 0 {
			switch randBool() {
			case duplicateVowel:
				wordB = append(wordB[:vI+1], wordB[vI:]...)
			case removeVowel:
				wordB = append(wordB[:vI], wordB[vI+1:]...)
			}
		}
	}
	return string(wordB)
}
