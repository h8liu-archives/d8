package main

import (
	"database/sql"
	"log"
	"os"

	_ "sqlite3"
)

func noError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	os.Remove("test.d8")
	db, e := sql.Open("sqlite3", "test.d8")
	noError(e)
	defer db.Close()

	x := func(q string) { _, e := db.Exec(q); noError(e) }

	x(`create table domains (id integer not null primary key, name text);`)
	x(`create table cnames (d integer, cname text);`)
	x(`create table ips (d integer, ip text, cname text);`)
	x(`create table servers (d integer, zone text, server text, ip text);`)
	x(`create table logs (d integer, log text);`)
}
