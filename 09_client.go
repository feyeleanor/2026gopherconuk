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
		if H.MatchString(string(b)) {
			s := H.FindAllStringIndex(string(b), FETCH_ALL_RESULTS)
			log.Printf("Found pattern %v times", len(s))

			for _, v := range s {
				log.Printf("Pattern matched: %s", b[v[0]:v[1]])
			}
		}
	})

	ProcessTargets(func(b ...byte) {
		ProcessText(H, b...)
	})
	os.Exit(E_OK)
}

func ProcessTargets(f func(...byte)) {
	ForArgs(func(fn string) {
		if W.MatchString(fn) {
			log.Printf("Load Web Page: %v", fn)
			c := http.Client{Timeout: time.Duration(TIMEOUT) * time.Second}
			ExitOnError(E_FILE_READ, func() (e error) {
				if r, e := c.Get(fn); e == nil {
					ReadStream(r.Body, func(b ...byte) {
						f(b...)
					})
				}
				return
			})

		} else {
			log.Printf("Load File: %v", fn)
			ExitOnError(E_FILE_READ, func() (e error) {
				if b, e := os.ReadFile(fn); e == nil {
					f(b...)
				}
				return
			})
		}
	})
	return
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

func ExitOnError(ec int, f func() error) {
	if e := f(); e != nil {
		log.Print(e)
		os.Exit(ec)
	}
}
