package main

import (
	"database/sql"
	"flag"
	"go/build"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type config struct {
	InputPath  string
	ResultPath string
	OutputPath string
	LogPath    string
	JobDB      string
}

var conf = &config{
	InputPath:  "input",
	ResultPath: "result",
	OutputPath: "output",
	LogPath:    "log",
	JobDB:      "jobs.db",
}

func ne(e error) {
	if e != nil {
		panic(e)
	}
}

func mkdir(s string) {
	stat, err := os.Stat(s)
	if err == nil {
		// already exists
		if stat.IsDir() {
			return
		}
		log.Fatalf("%q exists and it is not a directory", s)
	}
	if !os.IsNotExist(err) {
		log.Fatal(err)
	}

	e := os.Mkdir(s, 0700)
	if e != nil {
		log.Fatal(e)
	}
}

func dbinit() {
	dbpath := conf.JobDB
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
		id integer not null primary key, 
		name text, 
		state int, 
		total int,
		crawled int,
		birth text,
		death text
	);`)

	ne(db.Close())

	mkdir(conf.InputPath)
	mkdir(conf.ResultPath)
	mkdir(conf.OutputPath)
	mkdir(conf.LogPath)
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
