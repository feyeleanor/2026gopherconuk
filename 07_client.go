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
	ProcessTargets()
	os.Exit(E_OK)
}

func ProcessTargets() {
	if len(os.Args) > 0 {
		fn := os.Args[1:]

		if W.MatchString(fn[0]) {
			log.Printf("Load Web Page: %v", fn[0])
			c := http.Client{Timeout: time.Duration(TIMEOUT) * time.Second}
			if r, e := c.Get(fn[0]); e == nil {
				ReadStream(r.Body, func(b ...byte) {
					ProcessText(H, b...)
				})
			} else {
				log.Print(e)
				os.Exit(E_FILE_READ)
			}

		} else {
			log.Printf("Load File: %v", fn[0])
			if b, e := os.ReadFile(fn[0]); e == nil {
				ProcessText(H, b...)
			} else {
				log.Print(e)
				os.Exit(E_FILE_READ)
			}
		}
	}
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
