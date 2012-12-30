package main

import (
	"dico/dico"
	"log"
	"net/http"
	"os"
	"encoding/json"
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
		json.NewEncoder(w).Encode(d.Find(query))
	}))
	http.Handle("/", http.FileServer(http.Dir("resources")))

	const port = ":8123"
	log.Println("serving on http://localhost" + port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
