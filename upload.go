package main

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

func Upload(r *http.Request, w http.ResponseWriter) error {
	// Check if a file is being uploaded
	var file, header, err = r.FormFile("upload")
	if err != nil {
		return nil
	}
	// Extract path from url
	p, err := url.PathUnescape(r.URL.Path)
	if err != nil {
		return err
	}
	p = path.Clean(p)
	// Write file to disk
	f, err := os.OpenFile(path.Join(p, header.Filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	if _, err = io.Copy(f, file); err != nil {
		return err
	}
	// Write
	return f.Close()
}
