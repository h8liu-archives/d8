package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"d8/domain"
	"d8/tasks"
	"d8/term"
)

func noError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("d8> ")
		if !s.Scan() {
			break
		}

		line := s.Text()
		line = strings.TrimSpace(line)
		d, e := domain.Parse(line)
		if e != nil {
			fmt.Println("error: ", e)
			continue
		}

		term.T(tasks.NewInfo(d))
		fmt.Println()
	}

	noError(s.Err())

	fmt.Println()
}
