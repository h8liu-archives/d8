package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"d8/domain"
)

func LoadList(path string, eout io.Writer) ([]*domain.Domain, error) {
	fin, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer fin.Close()

	ret := make([]*domain.Domain, 0, 5000)

	lineno := 0
	s := bufio.NewScanner(fin)
	for s.Scan() {
		lineno++
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}

		d, e := domain.Parse(line)
		if e != nil {
			fmt.Fprintf(eout, "%s:%d: '%s': %v", path, lineno, line, e)
		} else {
			ret = append(ret, d)
		}
	}

	if s.Err() != nil {
		return ret, e
	}

	return ret, nil
}
