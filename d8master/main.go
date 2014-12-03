package main

import (
	"database/sql"
	"flag"
	"go/build"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func ne(e error) {
	if e != nil {
		panic(e)
	}
}

func dbinit() {
	dbpath := "jobs.db"
	db, err := sql.Open("sqlite3", dbpath)
	ne(err)

	q := func(sql string) {
		_, e := db.Exec(sql)
		if e != nil {
			log.Printf("sql: %s\n", sql)
			log.Fatal(e)
		}
	}

	q(`create table jobs (
		name text not null primary key,
		state int,
		total int,
		crawled int,
		sample text,
		error text,
		birth text,
		death text
	);`)

	ne(db.Close())
}

var (
	doInit    = flag.Bool("init", false, "perform db init")
	serveAddr = flag.String("http", ":8053", "the server address")
)

const (
	JobCreated = iota
	JobReady
	JobDone
)

func handleApi(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/")

	switch name {
	case "make":
		panic("todo")

	case "ls":
		panic("todo")

	default:
		w.WriteHeader(404)
	}
}

func wwwPath() string {
	pkg, e := build.Import("github.com/h8liu/d8/d8c", "", build.FindOnly)
	if e != nil {
		log.Fatal(e)
	}
	return filepath.Join(pkg.Dir, "www")
}

func main() {
	flag.Parse()

	if *doInit {
		dbinit()
		return
	}

	http.Handle("/", http.FileServer(http.Dir(wwwPath())))
	http.HandleFunc("/jobs/", handleApi)
	for {
		err := http.ListenAndServe(*serveAddr, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}
