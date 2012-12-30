package main

import (
	"github.com/oadam/longest-word/dico"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const maxLetters = 12

func main() {
	file, err := os.Open("mots.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var d = dico.New(file)

	http.Handle("/query", handler(d))
	http.Handle("/", http.FileServer(http.Dir("resources")))

	var port = ":" + os.Getenv("PORT")
	log.Println("serving on http://localhost" + port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

type handler dico.Dico

func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var d dico.Dico = dico.Dico(h);
	var receivedAt = time.Now()
	req.ParseForm()
	var query = req.Form.Get("q")
	if len(query) > maxLetters {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Too much (more than", maxLetters, ") letters")
		query = query[0:maxLetters] + "..."
	} else {
		var result = (&d).Find(query)
		json.NewEncoder(w).Encode(result)
	}
	log.Println("handled request \"", query, "\" in", time.Since(receivedAt))
}
