package main

import (
	"encoding/json"
	"fmt"
	"github.com/oadam/longest-word/dico"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const maxLetters = 1000002
const dicoFilename = "mots.txt"
const resourcesDirname = "resources"

var dic dico.Dico
var dicWait sync.WaitGroup

func main() {
	var startedAt = time.Now()

	dicoFile, errDico := os.Open(dicoFilename)
	if errDico != nil {
		panic(fmt.Sprintf("did not manage to open dico file at specified path \"%s\" (error: %s)", dicoFilename, errDico))
	}
	resourcesDir, errRes := os.Open(resourcesDirname)
	if errRes != nil {
		panic(fmt.Sprintf("did not find resources dir at path \"%s\" (error: %s)", resourcesDirname, errRes))
	}
	resourcesDir.Close()

	//init dico
	dicWait.Add(1)
	go func() {
		defer dicoFile.Close()
		dic = dico.New(dicoFile)
		dicWait.Add(-1)
		log.Printf("dictionnary ready (%s ellapsed)\n", time.Since(startedAt))
	}()

	http.Handle("/query", http.HandlerFunc(query))
	http.Handle("/", http.FileServer(http.Dir(resourcesDirname)))

	var port = ":" + os.Getenv("PORT")
	log.Printf("server started in %s, listening on %s\n", time.Since(startedAt), port)
	var err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func query(w http.ResponseWriter, req *http.Request) {
	var receivedAt = time.Now()
	req.ParseForm()
	var query = req.Form.Get("q")
	if len(query) > maxLetters {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Too much (more than", maxLetters, ") letters")
		query = query[0:maxLetters] + "..."
	} else {
		dicWait.Wait()
		var result = dic.Find(query)
		json.NewEncoder(w).Encode(result)
	}
	log.Println("handled request \"", query, "\" in", time.Since(receivedAt))
}
