package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const aliasSize int = 7

type DB struct {
	aliasKeysMap map[string]string
	fullURLKeysMap map[string]string
}

func initDB() *DB {
	db := new(DB)
	db.aliasKeysMap = make(map[string]string)
	db.fullURLKeysMap = make(map[string]string)
	return db
}

var db DB

func init() {
	db = *initDB()
}

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc(`/`, PostHandler)
	mux.HandleFunc(`/{id}`, GetHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

func RandomString(size int) string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	result := make([]rune, size)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range result {
		result[i] = chars[rnd.Intn(len(chars))]
	}
	return string(result)
}

func checkKeyIsExists(m map[string]string, k string) bool {
	_, ok := m[k]
	return ok
}

func PostHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	fullURL, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	if checkKeyIsExists(db.fullURLKeysMap, string(fullURL)) {
		fmt.Fprint(w, "http://"+req.Host+"/"+string(db.fullURLKeysMap[string(fullURL)]))
		return
	} else {
		alias := RandomString(aliasSize)
		db.fullURLKeysMap[string(fullURL)] = alias
		db.aliasKeysMap[string(alias)] = string(fullURL)
		fmt.Fprint(w, "http://"+req.Host+"/"+string(db.fullURLKeysMap[string(fullURL)]))
	}
}

func GetHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	alias, ok := mux.Vars(req)["id"];
	if !ok {
		http.Error(w, "Get query require Id", http.StatusBadRequest)
		return
	}
	if fullURL, ok := db.aliasKeysMap[alias]; ok {
		w.Header().Add("Location", fullURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	} else {
		fmt.Fprint(w, "URL by alias "+alias+" is not exists")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
