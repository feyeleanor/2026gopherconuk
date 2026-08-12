package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

const FETCH_ALL_RESULTS = -1
const HREF = `(i?)href[:blank:]*=[:blank:]*"[^"]*"`
const WEB_DOMAIN = `(i?)https?://.*/?`

const TIMEOUT = 3

const (
	E_OK = iota
	E_FILE_READ
)

var H, W *regexp.Regexp

func init() {
	H = regexp.MustCompile(HREF)
	W = regexp.MustCompile(WEB_DOMAIN)
}

func main() {
	ForArgs(func(fn string) {
		ProcessTarget(fn, ProcessText)
	})
	os.Exit(E_OK)
}

func ProcessTarget(fn string, f func(*regexp.Regexp, ...byte)) {
	if W.MatchString(fn) {
		log.Printf("Load Web Page: %v", fn)
		c := http.Client{Timeout: time.Duration(TIMEOUT) * time.Second}
		if r, e := c.Get(fn); e == nil {
			ReadStream(r.Body, func(b ...byte) {
				f(H, b...)
			})
		} else {
			log.Print(e)
			os.Exit(E_FILE_READ)
		}

	} else {
		log.Printf("Load File: %v", fn)
		if b, e := os.ReadFile(fn); e == nil {
			f(H, b...)
		} else {
			log.Print(e)
			os.Exit(E_FILE_READ)
		}
	}
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

func ForArgs(f func(string)) {
	if len(os.Args) > 0 {
		for _, fn := range os.Args[1:] {
			f(fn)
		}
	}
}

func ReadStream(in io.ReadCloser, f func(...byte)) (e error) {
	defer in.Close()
	b, e := io.ReadAll(in)
	if e == nil {
		f(b...)
	} else {
		log.Print(e)
	}
	return
}
