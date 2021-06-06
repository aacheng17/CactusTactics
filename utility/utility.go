package utility

import "strings"

const (
	MESSAGESEP = "\n"
	DELIM      = "\t"
	TAG        = "\v"
	BRTAG      = TAG + "br/" + TAG
	BTAG       = TAG + "b" + TAG
	ENDTAG     = TAG + "/" + TAG
)

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

func RemoveEscapes(s string) string {
	s = strings.ReplaceAll(s, DELIM, "")
	s = strings.ReplaceAll(s, TAG, "")
	s = strings.ReplaceAll(s, MESSAGESEP, BRTAG)
	return s
}
