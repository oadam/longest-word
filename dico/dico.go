package dico

import (
	"encoding/csv"
	"fmt"
	"github.com/cznic/sortutil"
	"github.com/fiam/gounidecode/unidecode"
	"io"
	"sort"
)

type node struct {
	words    []string
	children map[rune]*node
}

func (n *node) initChildren() {
	n.children = make(map[rune]*node)
}
func (n *node) addWord(w string) {
	n.words = append(n.words, w)
}
func (n *node) String() string {
	return fmt.Sprintln(n.words)
}

type Dico node

func (d *Dico) Find(letters string) []string {
	letters = unidecode.Unidecode(letters)
	var result []string
	var sorted = wordToSortedRunes(letters)
	var root = node(*d)
	var currents = []*node{&root}
	var seen = map[*node]bool{}
	for _, r := range sorted {
		var nexts = []*node{}
		for _, current := range currents {
			var next = current.children[r]
			if next != nil && !seen[next] {
				nexts = append(nexts, next)
				seen[next] = true
				for _, word := range next.words {
					result = append(result, word)
				}
			}
		}
		currents = append(currents, nexts...)
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

	var root *node = &node{}
	root.initChildren()
	for {
		words, e := reader.Read()
		if e != nil {
			break
		}
		var word = words[0]
		var sortedWord = wordToSortedRunes(word)
		var currentNode = root
		for _, r := range sortedWord {
			if currentNode.children == nil {
				currentNode.initChildren()
			}
			var childNode = currentNode.children[r]
			if childNode == nil {
				currentNode.children[r] = &node{}
				childNode = currentNode.children[r]
			}
			currentNode = childNode
		}
		currentNode.addWord(word)
	}
	return Dico(*root)
}
