package main

import (
	"os"
	"fmt"
	"bufio"
	"log"
	"strings"

	"d8/domain"
	"d8/term"
	"d8/tasks"
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
