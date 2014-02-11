package dico

import (
	"encoding/csv"
	"github.com/fiam/gounidecode/unidecode"
	"io"
	"sort"
)

type sortEntries []dicoEntry

func (result sortEntries) Len() int {
	return len(result)
}
func (result sortEntries) Less(i, j int) bool {
	if len(result[i].word) > len(result[j].word) {
		return true
	}
	return result[i].word < result[j].word
}
func (result sortEntries) Swap(i, j int) {
	result[i], result[j] = result[j], result[i]
}

type dicoEntry struct {
	word     string
	multiset map[rune]int
}
type Dico []dicoEntry

const maxResult = 100

func (d *Dico) Find(letters string) []string {
	var multiset = wordToMultiset(letters)
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

func wordToMultiset(word string) map[rune]int {
	var decoded = unidecode.Unidecode(word)
	var result = make(map[rune]int)
	for _, r := range decoded {
		result[r]++
	}
	return result
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
		entry.multiset = wordToMultiset(word)
		dico = append(dico, entry)
	}
	sort.Sort(sortEntries(dico))

	return dico
}
