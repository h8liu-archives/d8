package main

import (
	// "fmt"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"d8/client"
	. "d8/domain"
	//. "d8/packet/consts"
	"d8/tasks"
	"d8/term"
)

func noError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func ip(s string) net.IP {
	return net.ParseIP(s)
}

func loadList(path string) []*Domain {
	fin, e := os.Open(path)
	noError(e)

	ret := make([]*Domain, 0, 5000)
	s := bufio.NewScanner(fin)
	for s.Scan() {
		d := D(strings.TrimSpace(s.Text()))
		ret = append(ret, d)
	}
	noError(s.Err())
	fin.Close()

	return ret
}

func main() {
	list := loadList("list")

	c, e := client.New()
	noError(e)

	nquota := 30
	quotas := make(chan int, nquota)
	for i := 0; i < nquota; i++ {
		quotas <- i
	}

	for _, d := range list {
		i := <-quotas
		go crawl(d, c, quotas, i)
		time.Sleep(time.Second)
	}

	ids := make([]int, 0, nquota)
	for i := 0; i < nquota; i++ {
		ids = append(ids, <-quotas)
	}

	if len(ids) != nquota {
		panic("bug")
	}
}

func crawl(d *Domain, c *client.Client, q chan int, id int) {
	t := term.New(c)
	flog, e := os.Create(fmt.Sprintf("logs/%s", d))
	noError(e)
	t.Log = flog
	info := tasks.NewInfo(d)
	t.T(info)

	fmt.Fprintln(flog)

	flog.Close()
	fmt.Println(d)

	q <- id
}
