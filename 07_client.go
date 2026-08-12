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
	E_NOT_A_URL
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
	ForArgs(func(fn string) {
		if W.MatchString(fn) {
			log.Printf("Load Web Page: %v", fn)
			c := http.Client{Timeout: time.Duration(TIMEOUT) * time.Second}
			if r, e := c.Get(fn); e == nil {
				ReadStream(r.Body, func(b ...byte) {
					ProcessText(H, b...)
				})
			} else {
				log.Print(e)
				os.Exit(E_FILE_READ)
			}

		} else {
			os.Stderr.WriteString("cannot load local files\n")
			os.Stderr.WriteString("please run as: go run 05_client.go https://[url]\n")
			os.Exit(E_NOT_A_URL)
		}
	})
	return
}

func ProcessText(r *regexp.Regexp, b ...byte) {
	if r.MatchString(string(b)) {
		s := r.FindAllStringIndex(string(b), FETCH_ALL_RESULTS)
		log.Printf("Found pattern %v times", len(s))

		for _, v := range s {
			log.Printf("Pattern matched: %s", b[v[0]:v[1]])
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
