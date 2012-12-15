package dico

import (
	"os"
	"testing"
	"strings"
)

func TestRune(t *testing.T) {
	var truc = string(wordToSortedRunes("azertyuiopqsdfghjklmwxcvbn"))
	const expected = "abcdefghijklmnopqrstuvwxyz"
	if truc != expected {
		t.Errorf("expected %s but received %s", expected, truc)
	}
}
const exampleWords = "a\nab\naba\ntoto\n"

func TestNew(t *testing.T) {
	file := strings.NewReader(exampleWords)
	var d = New(file)
	var n = node(d)
	if _, ok := n.children['a']; !ok {
		t.Errorf("no node for a, map contains %v", len(n.children))
	}
	if words := n.children['a'].words; len(words) != 1 {
		t.Errorf("expected 1 word but found %v", words)
	}
}
func TestSimpleFind(t *testing.T) {
	file := strings.NewReader(exampleWords)
	var d = New(file)
	var results = d.Find("baa")
	if len(results) != 3 {
		t.Errorf("expected 3 results but received %v", results)
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
