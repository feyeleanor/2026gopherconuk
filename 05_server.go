package main

import (
	"io"
	"io/fs"
	"log"
	"net/http"
	"regexp"
)

const INDEX_FILE = `(i?)/$`
const SERVEABLE_FILE = `(i?)^.*\.html$`
const ADDRESS = ":3000"
const CURRENT_DIR = "."

//	https://eli.thegreenplace.net/2022/serving-static-files-and-web-apps-in-go/

var I, F *regexp.Regexp

func init() {
	I = regexp.MustCompile(INDEX_FILE)
	F = regexp.MustCompile(SERVEABLE_FILE)
}

func main() {
	log.Fatal(
		http.ListenAndServe(
			ADDRESS,
			http.FileServer(
				ProtectedFS(CURRENT_DIR))))
}

func ProtectedFS(root string) protectedFS {
	return protectedFS{http.Dir(root)}
}

func anyMatch(s string, p ...*regexp.Regexp) (r bool) {
	for _, v := range p {
		r = r || v.MatchString(s)
	}
	return
}

type protectedFile struct {
	http.File
}

func (f protectedFile) Readdir(n int) (fis []fs.FileInfo, e error) {
	files, e := f.File.Readdir(n)
	for _, file := range files {
		if anyMatch(file.Name(), F, I) {
			log.Printf("Match Succeed: Readdir(%v)", file.Name())
			fis = append(fis, file)
		} else {
			log.Printf("Match Failed: Readdir(%v)", file.Name())
		}
	}
	if e == nil && n > 0 && len(fis) == 0 {
		e = io.EOF
	}
	return
}

type protectedFS struct {
	http.FileSystem
}

func (fsys protectedFS) Open(name string) (http.File, error) {
	if !anyMatch(name, F, I) {
		log.Printf("Match Failed: Open(%v)", name)
		return nil, fs.ErrPermission
	}

	file, err := fsys.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	log.Printf("Match Succeed: Open(%v)", name)
	return protectedFile{file}, nil
}
