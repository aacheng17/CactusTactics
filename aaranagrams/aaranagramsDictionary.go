package aaranagrams

// declaring a struct
type AaranagramsDictionary struct {
	usedWords map[string]bool

	whattedWords map[int]string

	dictionary map[string]string

	minFreq int

	maxFreq int

	wordsLeft int
}

func (d *AaranagramsDictionary) generate(minWordLength int) {
	d.usedWords = make(map[string]bool)
	d.whattedWords = make(map[int]string)
	d.generateDictionary(minWordLength)
	d.wordsLeft = len(d.dictionary)
}

func (d *AaranagramsDictionary) generateDictionary(minWordLength int) {
	d.dictionary = make(map[string]string)
	for k, v := range dictionary {
		if len(k) >= minWordLength {
			d.dictionary[k] = v
		}
	}
}
