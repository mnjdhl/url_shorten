/*
Author: Manoj Dahal
*/

package main

import (
	"log"
	"net/http"
	"hash/fnv"
	"encoding/hex"
	"encoding/gob"
	"os"
	"time"
)

const SERVER_PORT = "9988"
const BASE_SHORT_URL = "ti.ny/"
var URL_STORE_LOC string
const SHORT_URL_END_POINT = "/shorturl"

type URLData struct {
	LongURL string
	ShortURL string
	DateCreated time.Time
}

var URLTable map[string] *URLData

/* Save URL Table into file */
func StoreURLTable() {

	df, err := os.OpenFile(URL_STORE_LOC, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Println("Failed to create/open the file "+URL_STORE_LOC+", ", err)
		return
	}

	defer df.Close()
	gob.Register(URLData{})
	enc := gob.NewEncoder(df)
	enc.Encode(URLTable)
}

/* Load URL Table from file */
func LoadURLTable() (urlTab map[string] *URLData) {

	df, err := os.OpenFile(URL_STORE_LOC, os.O_RDONLY, 0600)
	if err != nil {
		log.Println("Failed to open the file "+URL_STORE_LOC+", ", err)
		return nil
	}

	defer df.Close()

	gob.Register(URLData{})
	dec := gob.NewDecoder(df)
	err = dec.Decode(&urlTab)
	if err != nil {
		log.Println("Failed to decode from file "+URL_STORE_LOC+", ", err)
		return nil
	}

	i := 0
	for _, v := range urlTab {
		i++
		log.Println("Loaded URL Entry ", i, ": ", *v)
	}

	return urlTab
}


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
		StoreURLTable()
		return surl
	}

	log.Println("Found short URL in the table for ", lurl)
	return urlEntry.ShortURL
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

	_, err := os.Stat("/data")
	if os.IsNotExist(err) {
		URL_STORE_LOC = "./url_data.gob"
	} else {
		URL_STORE_LOC = "/data/url_data.gob"
	}

	URLTable = LoadURLTable()
	if URLTable == nil {
		URLTable = make(map[string]*URLData)
	}
	log.Println("Inialized..")
}

func main() {

	http.HandleFunc(SHORT_URL_END_POINT, HandleURLShortReqs)

	log.Println("Starting URL Shortening Service at port "+SERVER_PORT)
	log.Fatal(http.ListenAndServe(":"+SERVER_PORT, nil))
}
