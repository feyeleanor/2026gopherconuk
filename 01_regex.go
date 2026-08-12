package main

import (
	"io"
	"log"
	"regexp"
	"strings"
)

const REGEX = `(i?)href[:blank:]*=[:blank:]*"[^"]*"`

func main() {
	r := regexp.MustCompile(REGEX)
	log.Print(r)

	s := strings.NewReader(`
		<html>
			<head></head>
			<body>
				<a href="https://my.test.domain/some/path">
			</body>
		</html>`)

	f, e := io.ReadAll(s)
	if e != nil {
		panic(e)
	}

	fs := string(f)

	log.Print(fs)
	log.Printf("does string contain REGEX? %v", r.MatchString(fs))

	i := r.FindStringIndex(fs)
	log.Printf("Found REGEX between: %v and %v", i[0], i[1])
	log.Printf("Pattern matched: %s", f[i[0]:i[1]])
}
