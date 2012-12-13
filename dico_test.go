package dico

import (
	"fmt"
	"os"
	"testing"
)

func TestDico(t *testing.T) {
	file, err := os.Open("mots.txt")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	var d = New(file)
	var results = d.Find("reivilo")
	for result := range results {
		fmt.Println(result)
	}
}
