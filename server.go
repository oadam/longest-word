package main

import (
	"encoding/json"
	"fmt"
	"github.com/oadam/longest-word/dico"
	"log"
	"net/http"
	"os"
	"time"
)

const maxLetters = 12
const dicoFilename = "mots.txt"
const resourcesDirname = "resources"

func main() {
	dicoFile, errDico := os.Open(dicoFilename)
	if errDico != nil {
		panic(fmt.Sprintf("did not manage to open dico file at specified path \"%s\" (error: %s)", dicoFilename, errDico))
	}
	defer dicoFile.Close()
	if _, err := os.Open(resourcesDirname); err != nil {
		panic(fmt.Sprintf("did not find resources dir at path \"%s\" (error: %s)", resourcesDirname, err))
	}

	var d = dico.New(dicoFile)

	http.Handle("/query", handler(d))
	http.Handle("/", http.FileServer(http.Dir(resourcesDirname)))

	var port = ":" + os.Getenv("PORT")
	log.Println("serving on http://localhost" + port)
	var err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

type handler dico.Dico

func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var d dico.Dico = dico.Dico(h)
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
