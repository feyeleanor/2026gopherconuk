package main

import (
	"log"
	"os"
	"regexp"
)

const REGEX = `(i?)href[:blank:]*=[:blank:]*"[^"]*"`
const FETCH_ALL_RESULTS = -1

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
	s := string(b)
	if r.MatchString(s) {
		i := r.FindAllStringIndex(s, FETCH_ALL_RESULTS)
		log.Printf("Found pattern %v times", len(i))

		for _, v := range i {
			left, right := v[0], v[1]
			log.Printf("Found REGEX between: %v and %v", left, right)
			log.Printf("Pattern matched: %s", s[left:right])
		}
	}
}
