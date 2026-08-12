package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"
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
	ProcessTargets(func(n string, b ...byte) {
		ProcessText(H, n, b...)
	})
	os.Exit(E_OK)
}

func ProcessTargets(f func(string, ...byte)) {
	var wg sync.WaitGroup
	ForArgs(func(fn string) {
		wg.Go(func() {
			if W.MatchString(fn) {
				ForServer(fn, f)
			} else {
				ForFile(fn, f)
			}
		})
	})
	wg.Wait()
	return
}

func ForArgs(f func(string)) {
	if len(os.Args) > 0 {
		for _, fn := range os.Args[1:] {
			f(fn)
		}
	}
}

func ForServer(url string, f func(string, ...byte)) {
	log.Printf("Load Web Page: %v", url)
	c := http.Client{Timeout: time.Duration(TIMEOUT) * time.Second}
	ExitOnError(E_FILE_READ, func() (e error) {
		if r, e := c.Get(url); e == nil {
			ReadStream(r.Body, func(b ...byte) {
				f(url, b...)
			})
		}
		return
	})
}

func ForFile(n string, f func(string, ...byte)) {
	log.Printf("Load File: %v", n)
	ExitOnError(E_FILE_READ, func() (e error) {
		if b, e := os.ReadFile(n); e == nil {
			f(n, b...)
		}
		return
	})
}

func ProcessText(r *regexp.Regexp, n string, b ...byte) {
	s := string(b)
	if r.MatchString(s) {
		i := r.FindAllStringIndex(s, FETCH_ALL_RESULTS)
		log.Printf("%v: Found pattern %v times", n, len(i))

		for _, v := range i {
			left, right := v[0], v[1]
			log.Printf("%v: Found REGEX between: %v and %v", n, left, right)
			log.Printf("%v: Pattern matched: %s", n, s[left:right])
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
