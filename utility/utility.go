package utility

import (
	"fmt"
	"strings"
)

const (
	MESSAGESEP = "\n"
	DELIM      = "\t"
	TAG        = "\v"
	ENDTAG     = TAG + "/" + TAG
)

func Tag(tagString string) string {
	return TAG + tagString + TAG
}

func TagId(tagString string, tagId int) string {
	tagName := ""
	tagClasses := ""
	seenSpace := false
	for _, c := range tagString {
		if !seenSpace {
			if c == ' ' {
				seenSpace = true
			} else {
				tagName += string(c)
			}
		} else {
			tagClasses += string(c)
		}
	}
	if tagClasses != "" {
		tagClasses = " " + tagClasses
	}
	return fmt.Sprint(TAG, tagName, " ", tagId, tagClasses, TAG)
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func UrlIndexGetPath(url string, n int) string {
	if url[len(url)-1] != '/' {
		url += "/"
	}
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
	s = strings.ReplaceAll(s, MESSAGESEP, Tag("br"))
	return s
}

func ParseAndTag(s string) string {
	s = strings.Replace(s, "<i>", Tag("i"), -1)
	s = strings.Replace(s, "</i>", ENDTAG, -1)
	s = strings.Replace(s, "<b>", Tag("b"), -1)
	s = strings.Replace(s, "</b>", ENDTAG, -1)
	return s
}
