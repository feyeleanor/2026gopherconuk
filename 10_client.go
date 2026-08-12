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
		ProcessTarget(fn, FindLinks)
	})
	os.Exit(E_OK)
}

func ProcessTarget(fn string, f func(...byte)) {
	if W.MatchString(fn) {
		ForServer(fn, f)
	} else {
		ForFile(fn, f)
	}
}

func ForArgs(f func(string)) {
	if len(os.Args) > 0 {
		for _, fn := range os.Args[1:] {
			f(fn)
		}
	}
}

func ForServer(url string, f func(...byte)) {
	log.Printf("Load Web Page: %v", url)
	c := http.Client{Timeout: time.Duration(TIMEOUT) * time.Second}
	ExitOnError(E_FILE_READ, func() (e error) {
		if r, e := c.Get(url); e == nil {
			ReadStream(r.Body, func(b ...byte) {
				f(b...)
			})
		}
		return
	})
}

func ForFile(n string, f func(...byte)) {
	log.Printf("Load File: %v", n)
	ExitOnError(E_FILE_READ, func() (e error) {
		if r, e := os.Open(n); e == nil {
			ReadStream(r, func(b ...byte) {
				f(b...)
			})
		}
		return
	})
}

func FindLinks(b ...byte) {
	s := string(b)
	i := H.FindAllStringIndex(s, FETCH_ALL_RESULTS)
	log.Printf("Found pattern %v times", len(i))

	for _, v := range i {
		left, right := v[0], v[1]
		log.Printf("Found REGEX between: %v and %v", left, right)
		log.Printf("Pattern matched: %s", s[left:right])
	}
}

func ReadStream(in io.ReadCloser, f func(...byte)) (e error) {
	defer in.Close()
	if b, e := io.ReadAll(in); e == nil {
		f(b...)
	} else {
		log.Print(e)
	}
	return
}

func ExitOnError(ec int, f func() error) {
	if e := f(); e != nil {
		log.Print(e)
		os.Exit(ec)
	}
}
