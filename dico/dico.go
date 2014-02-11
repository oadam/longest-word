package dico

import (
	"bytes"
	"encoding/csv"
	"github.com/cznic/sortutil"
	"github.com/fiam/gounidecode/unidecode"
	"io"
	"sort"
)

type node map[string][]string

type Dico map[string][]string

func (d *Dico) Find(letters string) []string {
	letters = unidecode.Unidecode(letters)
	var sorted = wordToSortedRunes(letters)
	var asMap = map[string][]string(*d)
	var result []string
	var b bytes.Buffer
	var l = uint(len(sorted))
	var subsets = make(map[string]bool)
	for i := 0; i < 1<<l; i++ {
		b.Reset()
		for j, r := range sorted {
			var uj = uint(j)
			if i&(1<<uj) != 0 {
				b.WriteRune(r)
			}
		}
		var subset = b.String()
		subsets[subset] = true
	}
	for subset, _ := range subsets {
		result = append(result, asMap[subset]...)
	}

	sort.Strings(result)
	return result
}

func wordToSortedRunes(word string) []rune {
	var decoded = unidecode.Unidecode(word)
	var sortedWord = make([]rune, len(decoded))
	copy(sortedWord, []rune(decoded))
	sort.Sort(sortutil.RuneSlice(sortedWord))
	return sortedWord
}

func New(file io.Reader) Dico {
	var reader = csv.NewReader(file)
	reader.FieldsPerRecord = 1

	var result = make(map[string][]string)
	for {
		words, e := reader.Read()
		if e != nil {
			break
		}
		var word = words[0]
		var sortedWord = wordToSortedRunes(word)
		var sortedString = string(sortedWord)
		result[sortedString] = append(result[sortedString], word)
	}
	return Dico(result)
}
