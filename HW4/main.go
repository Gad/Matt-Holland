package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type Item string   // define as new type
type Price float64 // define as new type


func (p Price) String() string{
	return fmt.Sprintf("%.2f€", p)
}
type Db map[Item]Price

var db = Db{"car": 10000, "computer": 500}

func (db *Db) Create(i Item, p Price) {
	(*db)[i] = p

}

func (db *Db) Read(i Item) (Price, error) {

	p, exists := (*db)[i]
	if !exists {
		return 0, fmt.Errorf("Error reading Price of %s : no such Item in database", i)
	}
	return p, nil
}

func (db *Db) Update(i Item, p Price) error {

	if _, exists := (*db)[i]; exists {
		(*db)[i] = p
		return nil
	}
	return fmt.Errorf("Error updating %s with Price %s : no such Item in database", i, p)
}

func (db *Db) Delete(i Item) {
	delete(*db, i)
}

func parseQuery(q url.Values) (Item, Price, error) {
	var (
		i Item
		p Price
	)
	for k, v := range q {
		switch k {
		case "item":
			i = Item(v[0])
		case "price":
			price, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				return i, p, fmt.Errorf("Price Not a number")
			}
			p = Price(price)
		default:
			log.Printf("unsupported query parameter %s\n", k)
		}
	}
	return i, p, nil
}

var update = func(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	i, p, err := parseQuery(q)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if err := db.Update(i, p); err != nil {
		log.Println(err.Error())
		return
	}
	log.Printf("Updated %s with Price %s", i, p)
}

var read = func(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	i, _, err := parseQuery(q)
	if err != nil {
		log.Println(err.Error())
		return
	}
	p, err := db.Read(i)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Printf("Price for Item %s is %s", i, p)
}

var create = func(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	i, p, err := parseQuery(q)
	if i == "" {
		log.Println("Cannot create empty Item")
		return
	}
	if err != nil {
		log.Println(err.Error())
		return
	}
	if _, err := db.Read(i); err == nil {
		log.Printf("Item %s already in db, use update instead", i)
		return
	}
	db.Create(i, p)
	log.Printf("Created Item %s with Price %s", i, p)
}

var deleteKey = func(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	i, _, err := parseQuery(q)
	if i == "" {
		log.Println("Item name required with /delete query")
		return
	}
	if err != nil {
		log.Println(err.Error())
		return
	}
	db.Delete(i)
	log.Printf("Deleted Item %s", i)
}

func withLogging(h http.HandlerFunc) http.HandlerFunc {
	logFunc := func(w http.ResponseWriter, r *http.Request){
		mw := io.MultiWriter(os.Stdout, w)
		log.SetOutput(mw)
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(logFunc)
}

func main() {

	http.HandleFunc("/create", withLogging(create))
	http.HandleFunc("/read", withLogging(read))
	http.HandleFunc("/update", withLogging(update))
	http.HandleFunc("/delete", withLogging(deleteKey)) // delete already in use by Go


	log.Fatal(http.ListenAndServe(":8080", nil))

}
