package main

import "strings"

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func urlIndexGetPath(url string, n int) string {
	firstSlashIndex := strings.Index(url, "/")
	url = url[firstSlashIndex+1:]
	for i := 0; i < n; i++ {
		slashIndex := strings.Index(url, "/")
		if slashIndex == -1 {
			break
		}
		url = url[slashIndex+1:]
	}
	lastSlashIndex := strings.Index(url, "/")
	if lastSlashIndex != -1 {
		url = url[:lastSlashIndex]
	}
	return url
}

func makeRange(min, max int) []int {
	a := make([]int, max-min)
	for i := range a {
		a[i] = min + i
	}
	return a
}
