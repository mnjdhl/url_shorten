/*
Author: Manoj Dahal
*/

package main

import (
	"log"
	"net/http"
	"hash/fnv"
	"encoding/hex"
	"time"
)

const SERVER_PORT = "9988"
const BASE_SHORT_URL = "ti.ny/"

type URLData struct {
	longURL string
	shortURL string
	dateCreated time.Time
}

var URLTable map[string] *URLData

/* Generate non-cryptographic hash for a string */
func GenerateHash(longURL string) string {

	//hfn := fnv.New64()
	hfn := fnv.New32()

	hfn.Write([]byte(longURL))
	sm := hfn.Sum(nil)

	log.Println("Hash generated for ", longURL)
	return hex.EncodeToString(sm)
}

/*Get/Create short url for a given long url*/
func GetShortURL(lurl string) string {

	urlEntry := URLTable[lurl]
	if urlEntry == nil {
		surl := BASE_SHORT_URL+GenerateHash(lurl)
		newUrlEntry := URLData{lurl, surl, time.Now()}
		URLTable[lurl] = &newUrlEntry
		return surl
	}

	log.Println("Found short URL in the table for ", lurl)
	return urlEntry.shortURL
	//return ""
}

/*HTTP Handler function for url shorting request*/
func HandleURLShortReqs(hrw http.ResponseWriter, hreq *http.Request) {

	if hreq.Method == "GET" {
		if hreq.URL.Query()["longURL"] == nil {
			http.Error(hrw, "Invalid Long URL", http.StatusBadRequest)
			return
		}

		longURL := hreq.URL.Query().Get("longURL")
		shortURL := GetShortURL(longURL)
		hrw.WriteHeader(http.StatusOK)
		hrw.Header().Set("Content-Type", "text/plain")
		hrw.Write([]byte(shortURL))
	}

}

/*Init function*/
func init() {

	URLTable = make(map[string]*URLData)
	log.Println("Inialized..")
}

func main() {

	http.HandleFunc("/shorturl", HandleURLShortReqs)

	log.Println("Starting URL Shortening Service at port "+SERVER_PORT)
	http.ListenAndServe(":"+SERVER_PORT, nil)


}
