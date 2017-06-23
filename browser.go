package main

import (
	"fmt"
	"html"
	"net/http"
	"net/url"
	"os"
	"path"
)

func Browser(w http.ResponseWriter, r *http.Request) {
	var p, err = url.PathUnescape(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p = path.Clean(p)
	info, err := os.Stat(p)
	// Check for non-existent files and stat errors
	if os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Serve file contents if not a directory
	if !info.IsDir() {
		http.ServeFile(w, r, p)
		return
	}
	// Check if user asked for a deletion
	if delQuery := r.URL.Query()["delete"]; len(delQuery) > 0 {
		var delName, err = url.QueryUnescape(delQuery[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		os.RemoveAll(path.Join(p, delName))
	}
	// List directory contents
	dir, err := os.Open(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	contents, err := dir.Readdir(0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "<table>")
	fmt.Fprintf(w, "<td><a href=\"/%s\">..</a></td>\n", url.PathEscape(path.Join(p, "..")))
	for _, el := range contents {
		fmt.Fprintln(w, "<tr>")
		fmt.Fprintf(w, "<td><a href=\"/%s\">%s</a></td>\n", url.PathEscape(path.Join(p, el.Name())), html.EscapeString(el.Name()))
		if el.IsDir() {
			fmt.Fprintf(w, "<td>dir</td>\n")
		} else {
			fmt.Fprintf(w, "<td>%d</td>\n", el.Size())
		}
		fmt.Fprintf(w, "<td>%s</td>\n", el.ModTime().Format("2006-01-02 15:04:05"))
		fmt.Fprintf(w, "<td><a href=\"/%s?delete=%s\">delete</a></td>\n", url.PathEscape(p), url.QueryEscape(el.Name()))
		fmt.Fprintln(w, "</tr>")
	}
	fmt.Fprintln(w, "</table>")
}
