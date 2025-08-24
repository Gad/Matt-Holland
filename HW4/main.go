package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type item = string   // alias
type price = float64 //alias

type Db map[item]price

var db = Db{"car": 10000, "computer": 500}

func (db *Db) Create(i item, p price) {
	(*db)[i] = p

}

func (db *Db) Read(i item) (float64, error) {

	p, exists := (*db)[i]
	if !exists {
		return 0, fmt.Errorf("Error reading price of %s : no such item in database", i)
	}
	return p, nil
}

func (db *Db) Update(i item, p price) error {

	if _, exists := (*db)[i]; exists {
		(*db)[i] = p
		return nil
	}
	return fmt.Errorf("Error updating %s with price %f : no such item in database", i, p)
}

func (db *Db) Delete(i item) {
	delete(*db, i)
}

func parseQuery(q url.Values) (string, float64, error) {
	var (
		i item
		p price
	)
	for k, v := range q {
		switch k {
		case "item":
			i = v[0]
		case "price":
			price, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				return i, p, fmt.Errorf("Price Not a number")
			}
			p = price
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
	log.Printf("Updated %s with price %f", i, p)
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
	log.Printf("Price for item %s is %f", i, p)
}

var create = func(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	i, p, err := parseQuery(q)
	if i == "" {
		log.Println("Cannot create empty item")
		return
	}
	if err != nil {
		log.Println(err.Error())
		return
	}
	if _, err := db.Read(i); err == nil {
		log.Printf("item %s already in db, use update instead", i)
		return
	}
	db.Create(i, p)
	log.Printf("Created item %s with price %f", i, p)
}

var deleteKey = func(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	i, _, err := parseQuery(q)
	if i == "" {
		log.Println("item name required with /delete query")
		return
	}
	if err != nil {
		log.Println(err.Error())
		return
	}
	db.Delete(i)
	log.Printf("Deleted item %s", i)

}

func main() {

	http.HandleFunc("/create", create)
	http.HandleFunc("/read", read)
	http.HandleFunc("/update", update)
	http.HandleFunc("/delete", deleteKey) // delete already in use by Go

	log.Fatal(http.ListenAndServe(":8080", nil))

}
