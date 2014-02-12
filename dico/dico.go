package dico

import (
	"encoding/csv"
	"github.com/fiam/gounidecode/unidecode"
	"io"
	"sort"
	"strings"
)

type byWordLength []dicoEntry

func (result byWordLength) Len() int {
	return len(result)
}
func (result byWordLength) Less(i, j int) bool {
	if result[i].length != result[j].length {
		return result[i].length > result[j].length
	} else {
		return result[i].word < result[j].word
	}
}
func (result byWordLength) Swap(i, j int) {
	result[i], result[j] = result[j], result[i]
}

type dicoEntry struct {
	word     string
	length   int
	multiset map[rune]int
}
type Dico []dicoEntry

const maxResult = 100

func (d *Dico) Find(letters string) []string {
	var multiset, _ = wordToMultiset(letters)
	var result []string
	for _, entry := range []dicoEntry(*d) {
		var eSet = entry.multiset
		var ok = true
		for r, nb := range eSet {
			if multiset[r] < nb {
				ok = false
				break
			}
		}
		if ok {
			result = append(result, entry.word)
			if len(result) > maxResult {
				break
			}
		}
	}
	return result
}

func wordToMultiset(word string) (map[rune]int, int) {
	var decoded = unidecode.Unidecode(word)
	decoded = strings.ToLower(decoded)
	var result = make(map[rune]int)
	var l = 0
	for _, r := range decoded {
		result[r]++
		l++
	}
	return result, l
}

func New(file io.Reader) Dico {
	var reader = csv.NewReader(file)
	reader.FieldsPerRecord = 1

	var dico Dico
	for {
		words, e := reader.Read()
		if e != nil {
			break
		}
		var word = words[0]
		var entry dicoEntry
		entry.word = word
		entry.multiset, entry.length = wordToMultiset(word)
		dico = append(dico, entry)
	}
	sort.Sort(byWordLength(dico))
	return dico
}
