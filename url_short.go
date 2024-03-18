/*
Author: Manoj Dahal
*/

package main

import (
	"encoding/gob"
	"encoding/hex"
	"hash/fnv"
	"log"
	"net/http"
	"os"
	"time"
)

const (
    HASH_SIZE32 = iota
    HASH_SIZE64
)
const SERVER_PORT = "9988"
const BASE_SHORT_URL = "ti.ny/"

var URL_STORE_LOC string

const SHORT_URL_END_POINT = "/shorturl"

var updateUrlChan chan bool
var updateUrlCounter int
var updateUrlThreshold int

type URLData struct {
	LongURL     string
	ShortURL    string
	DateCreated time.Time
}

var URLTable map[string]*URLData

func StoreURLRoutine() {

	for {
		updateFlag := <-updateUrlChan
		if updateFlag {
			StoreURLTable()
		}
	}
}

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
	log.Println("Updated URL Entries in " + URL_STORE_LOC)
}

/* Load URL Table from file */
func LoadURLTable() (urlTab map[string]*URLData) {

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
func GenerateHash(longURL string, size uint32) string {
    var sm []byte

    if size == HASH_SIZE32 {
	    hfn := fnv.New32()
	    hfn.Write([]byte(longURL))
	    sm = hfn.Sum(nil)
    } else if size == HASH_SIZE64 {
	    hfn := fnv.New64()
	    hfn.Write([]byte(longURL))
	    sm = hfn.Sum(nil)
    }


	log.Println("Hash generated for ", longURL)
	return hex.EncodeToString(sm)
}

/*Get/Create short url for a given long url*/
func GetShortURL(lurl string) string {

	urlEntry := URLTable[lurl]
	if urlEntry == nil {
		uhash := GenerateHash(lurl, HASH_SIZE32)
		surl := BASE_SHORT_URL + uhash
		newUrlEntry := URLData{lurl, surl, time.Now()}
		URLTable[lurl] = &newUrlEntry
		URLTable[uhash] = &newUrlEntry
		updateUrlCounter++
		if (updateUrlCounter % updateUrlThreshold) == 0 {
			updateUrlChan <- true
		}
		return surl
	}

	log.Println("Found short URL in the table for ", lurl)
	return urlEntry.ShortURL
}

func GetLongURL(surl string) string {

	surl_ex := BASE_SHORT_URL + surl
	urlEntry := URLTable[surl]
	if urlEntry != nil {
		if urlEntry.ShortURL == surl_ex {
			log.Println("Found the long url for the short url ", surl)
			return urlEntry.LongURL
		}
	}

	for _, v := range URLTable {
		log.Println(v)
		if v.ShortURL == surl_ex {
			log.Println("Found long URL in the table for ", surl)
			return v.LongURL
		}
	}
	return ""
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

func RedirectShortURL(hrw http.ResponseWriter, hreq *http.Request) {

	if hreq.Method == "GET" {
		shortURL := hreq.URL.Path[1:]
		longURL := GetLongURL(shortURL)
		if longURL != "" {
			log.Println("Redirecting ", shortURL, " to ", longURL)
			http.Redirect(hrw, hreq, longURL, http.StatusSeeOther)
		} else {
			log.Println("Redirecting failed as ", shortURL, " not found in the table")
			http.Error(hrw, "Invalid Short URL", http.StatusBadRequest)
		}
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
	updateUrlChan = make(chan bool, 1)
	updateUrlCounter = 0
	updateUrlThreshold = 10
	go StoreURLRoutine()
	log.Println("Inialized..")

}

func main() {

	http.HandleFunc(SHORT_URL_END_POINT, HandleURLShortReqs)
	http.HandleFunc("/", RedirectShortURL)

	log.Println("Starting URL Shortening Service at port " + SERVER_PORT)
	log.Fatal(http.ListenAndServe(":"+SERVER_PORT, nil))
}
