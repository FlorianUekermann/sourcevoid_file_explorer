package main

import (
	"os"
	"sort"
)

type FileInfoListSortByModTime []os.FileInfo

func (l FileInfoListSortByModTime) Sort() {
	sort.Sort(l)
}

func (l FileInfoListSortByModTime) Less(i, j int) bool {
	return l[i].ModTime().After(l[j].ModTime())
}

func (l FileInfoListSortByModTime) Len() int {
	return len(l)
}

func (l FileInfoListSortByModTime) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
