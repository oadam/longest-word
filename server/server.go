package main

import (
	"fmt"
	"log"
	"os"
	"dico/dico"
	"net/http"
)

func main() {
	file, err := os.Open("mots.txt")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	var d = dico.New(file)

	http.Handle("/query", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		var query = req.Form.Get("q")
		fmt.Fprintf(w, "%v\n", d.Find(query))
	}))

	const port = ":8123"
	log.Println("serving on http://localhost" + port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
