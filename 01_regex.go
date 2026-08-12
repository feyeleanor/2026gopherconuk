package main

import (
	"io"
	"log"
	"regexp"
	"strings"
)

const REGEX = `(i?)href[:blank:]*=[:blank:]*"[^"]*"`
const TEST_PAGE = `
		<html>
			<head></head>
			<body>
				<a href="https://my.test.domain/some/path">
			</body>
		</html>`

func main() {
	r := regexp.MustCompile(REGEX)
	log.Print(r)

	s := strings.NewReader(TEST_PAGE)
	f, _ := io.ReadAll(s)

	fs := string(f)
	log.Print(fs)
	log.Printf("does string contain REGEX? %v", r.MatchString(fs))

	i := r.FindStringIndex(fs)
	left, right := i[0], i[1]
	log.Printf("Found REGEX between: %v and %v", left, right)
	log.Printf("Pattern matched: %s", f[left:right])
}
