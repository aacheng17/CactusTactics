package utility

import "strings"

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func UrlIndexGetPath(url string, n int) string {
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

func MakeRange(min, max int) []int {
	a := make([]int, max-min)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func EscapeString(s string) string {
	htmlEscaper := strings.NewReplacer(
		`&`, "&amp;",
		`'`, "&#39;",
		`<`, "&lt;",
		`>`, "&gt;",
		`"`, "&#34;",
	)
	return htmlEscaper.Replace(s)
}

func UnEscapeString(s string) string {
	htmlUnEscaper := strings.NewReplacer(
		`&amp;`, "&",
		`&#39;`, "'",
		`&lt;`, "<",
		`&gt;`, ">",
		`&#34;`, "\"",
	)
	return htmlUnEscaper.Replace(s)
}
