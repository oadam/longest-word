package main

import (
	"dico/dico"
	"encoding/json"
	"log"
	"time"
	"net/http"
	"fmt"
	"os"
)
const maxLetters = 12

func main() {
	file, err := os.Open("mots.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var d = dico.New(file)

	http.Handle("/query", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var receivedAt = time.Now()
		req.ParseForm()
		var query = req.Form.Get("q")
		if len(query) > maxLetters {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Too much (more than", maxLetters,") letters")
			query = query[0:maxLetters] + "..."
		} else {
			json.NewEncoder(w).Encode(d.Find(query))
		}
		log.Println("handled request \"",query,"\" in", time.Since(receivedAt))
	}))
	http.Handle("/", http.FileServer(http.Dir("resources")))

	var port = ":" + os.Getenv("PORT")
	log.Println("serving on http://localhost" + port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
