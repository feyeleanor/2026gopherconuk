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

var W *regexp.Regexp

func init() {
	W = regexp.MustCompile(WEB_DOMAIN)
}

func main() {
	H := regexp.MustCompile(HREF)

	ProcessTargets(func(b ...byte) {
		ProcessText(H, b...)
	})
	os.Exit(E_OK)
}

func ProcessTargets(f func(...byte)) {
	if len(os.Args) > 0 {
		for _, fn := range os.Args[1:] {
			go func(n string) {
				if W.MatchString(n) {
					ForServer(n, f)
				} else {
					ForFile(n, f)
				}
			}(fn)
		}
	}
	return
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
		if b, e := os.ReadFile(n); e == nil {
			f(b...)
		}
		return
	})
}

func ProcessText(r *regexp.Regexp, b ...byte) {
	if r.MatchString(string(b)) {
		s := r.FindAllStringIndex(string(b), FETCH_ALL_RESULTS)
		log.Printf("Found pattern %v times", len(s))

		for _, v := range s {
			log.Printf("Pattern matched: %v", string(b[v[0]:v[1]]))
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

func ExitOnError(ec int, f func() error) {
	if e := f(); e != nil {
		log.Print(e)
		os.Exit(ec)
	}
}
