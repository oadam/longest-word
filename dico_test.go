package dico

import (
	"testing"
	"os"
	"fmt"
)

func TestDico(t *testing.T) {
	file, err := os.Open("mots.txt")
        defer file.Close()
        if err != nil {
                panic(err)
        }
        var d = New(file)
        var results = d.find("reivilo")
        fmt.Println(results)
}

