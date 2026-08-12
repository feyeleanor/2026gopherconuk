package main

import (
	"log"
	"os"
	"regexp"
)

const REGEX = `(i?)href[:blank:]*=[:blank:]*"[^"]*"`

const (
	E_OK = iota
	E_FILE_READ
)

func main() {
	//	compile Regex
	r := regexp.MustCompile(REGEX)
	log.Print(r)

	//	check for a filename as a parameter
	if len(os.Args) > 0 {
		//	get the filename
		fn := os.Args[1]

		//	load the file
		if b, e := os.ReadFile(fn); e == nil {
			ProcessText(r, b...)
		} else {
			os.Exit(E_FILE_READ)
		}
	}

	os.Exit(E_OK)
}

func ProcessText(r *regexp.Regexp, b ...byte) {
	if r.MatchString(string(b)) {
		s := string(b)
		i := r.FindStringIndex(s)
		left, right := i[0], i[1]
		log.Printf("Found REGEX between: %v and %v", left, right)
		log.Printf("Pattern matched: %s", s[left:right])
	}
}
