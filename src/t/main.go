package main

import (
	"log"

	_ "sqlite3"
)

func noError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	db, e := sql.Open("sqlite3", "./test.d8")
	noError(e)

	_, e := db.Exec(`
		create table domains (
			id integer not null primary key,
			domain text
		);

		create table cnames (
			domain integer,
			cname text
		);
		
		create table ips (
			domain integer,
			ip text,
			cname text
		);

		create table servers (
			domain integer,
			zone text,
			server text,
			ip text
		)
		`)

	noError(e)
