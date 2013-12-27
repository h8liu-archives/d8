package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "sqlite3"
)

func noError(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func test() {
	os.Remove("test.d8")
	db, e := sql.Open("sqlite3", "test.d8")
	noError(e)
	defer db.Close()
	x := func(s string) { _, e := db.Exec(s); noError(e) }
	/*
		q := func(s string) *sql.Rows {
			r, e := db.Query(s); noError(e)
			return r
		}
	*/

	x(`create table domains (id integer not null primary key, name text)`)
	x(`create table cnames (d integer, cname text)`)
	x(`create table ips (d integer, ip text, cname text)`)
	x(`create table servers (d integer, zone text, server text, ip text)`)
	x(`create table logs (d integer, log text)`)
}

func main() {
	os.Remove("test.d8")
	db, e := sql.Open("sqlite3", "test.d8")
	noError(e)
	defer db.Close()
	x := func(s string) { _, e := db.Exec(s); noError(e) }
	q := func(s string) *sql.Rows {
		r, e := db.Query(s)
		noError(e)
		return r
	}

	x(`create table test (id integer not null primary key, t text)`)
	r := q(`insert into test (t) values 
		('row1'),
		('row2'),
		('row3')
	`)

	var id int
	for r.Next() {
		e = r.Scan(&id)
		noError(e)
		fmt.Println("hello")
		fmt.Println(id)
	}
	noError(r.Err())
}
