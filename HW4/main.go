package main

import "fmt"


type item = string // alias
type price = float64 //alias

type Db map[item]price 
var db = Db{"car":10000,"computer":500}

func(db *Db) Create(i item, p price) error{
	return nil
}

func (db *Db) Read(i item) (float64, error) {
	return 0,nil
}

func(db *Db) Update(i item, p price) error{
	return nil
}

func(db *Db) Delete(i item) error{
	return nil
}



func main(){

	fmt.Printf("%T %[1]v", db)


}
