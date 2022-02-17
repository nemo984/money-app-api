package util

import "net/http"

//https://stackoverflow.com/a/51170557
type Fs struct {
	http.Dir
}

var notFoundFile, notFoundErr = http.Dir("dummy").Open("does-not-exist")

func (m Fs) Open(name string) (result http.File, err error) {
	f, err := m.Dir.Open(name)
	if err != nil {
		return
	}

	fi, err := f.Stat()
	if err != nil {
		return
	}
	if fi.IsDir() {
		// Return a response that would have been if directory would not exist:
		return notFoundFile, notFoundErr
	}
	return f, nil
}
