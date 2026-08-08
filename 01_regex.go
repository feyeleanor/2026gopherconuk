package main

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

const REGEX = `(i?)href[:blank:]*=[:blank:]*"[^"]*"`

func main() {
	r := regexp.MustCompile(REGEX)
	fmt.Println(r)

	s := strings.NewReader(`a href="https://my.test.domain/some/path"`)

	f, e := io.ReadAll(s)
	if e != nil {
		panic(e)
	}

	fmt.Println(f)
}