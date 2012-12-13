package dico

import (
	"encoding/csv"
	"io"
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

func (d *Dico) Find(letters string) chan string {
	var resultChan = make(chan string)
	go func() {
		var sorted = wordToSortedRunes(letters)
		var root = node(*d)
		var currents = []*node{&root}
		for _, r := range sorted {
			var nexts = []*node{}
			for _, current := range currents {
				var next, ok = current.children[r]
				if ok {
					nexts = append(nexts, &next)
					for _, word := range next.words {
						resultChan <- word
					}
				}
			}
			currents = append(currents, nexts...)
		}
		close(resultChan)
	}()
	return resultChan
}

func wordToSortedRunes(word string) []rune {
	var sortedWord = make([]rune, len(word))
	copy(sortedWord, []rune(word))
	sort.Sort(sortRunes(sortedWord))
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
			var childNode, ok = currentNode.children[r]
			if !ok {
				currentNode.children[r] = node{}
				childNode = currentNode.children[r]
			}
			currentNode = &childNode
		}
		currentNode.addWord(word)
	}
	return Dico(*root)
}
