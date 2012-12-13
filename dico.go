package dico

import (
	"io"
	"encoding/csv"
	"sort"
)

type node struct {
	words    []string
	children map[rune]node
}

func (n *node) initChildren() {
	n.children = make(map[rune]node)
}
func (n *node) addWord(w string) {
	n.words = append(n.words, w)
}

type sortRunes []rune

func (runes sortRunes) Len() int {
	return len(runes)
}
func (runes sortRunes) Less(i, j int) bool {
	return runes[i] < runes[j]
}
func (runes sortRunes) Swap(i, j int) {
	runes[i], runes[j] = runes[j], runes[i]
}

type Dico node

func New(file io.Reader) Dico {
	var reader = csv.NewReader(file)
	reader.FieldsPerRecord = 1

	var root node
	root.initChildren()
	for {
		words, e := reader.Read()
		if e != nil {
			break
		}
		var word = words[0]
		var sortedWord = make([]rune, len(word))
		copy(sortedWord, []rune(word))
		sort.Sort(sortRunes(sortedWord))
		var currentNode = root
		for i := 0; i < len(sortedWord); i++ {
			var r = sortedWord[i]
			if currentNode.children == nil {
				currentNode.initChildren()
			}
			currentNode = currentNode.children[r]
		}
		currentNode.addWord(word)
	}
	return Dico(root)
}

