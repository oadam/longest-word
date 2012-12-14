package dico

import (
	"fmt"
	"os"
	"testing"
)

func TestRune(t *testing.T) {
	var truc = wordToSortedRunes("azertyuiopqsdfghjklmwxcvbn")
	fmt.Println(string(truc))
}

func TestNew(t *testing.T) {
	file, err := os.Open("test_mots.txt")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	var d = New(file)
	var n = node(d)
	if _, ok := n.children['a']; !ok {
		fmt.Printf("no node for a, map contains %v", len(n.children))
	}
	if words := n.children['a'].words; len(words) != 1 {
		fmt.Printf("expected 1 word but found %v", words)
	}
}
func TestFind(t *testing.T) {
	file, err := os.Open("test_mots.txt")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	var d = New(file)
	var results = d.Find("baaissables")
	for result := range results {
		fmt.Println(result)
	}
}
