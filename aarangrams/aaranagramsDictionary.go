package aaranagrams

// declaring a struct
type AaranagramsDictionary struct {
	usedWords map[string]bool

	whattedWords map[int]string

	dictionary map[string]string

	letters map[string]int

	minFreq int

	maxFreq int

	wordsLeft int
}

func (d *AaranagramsDictionary) generate(minWordLength int) {
	d.usedWords = make(map[string]bool)
	d.whattedWords = make(map[int]string)
	d.generateDictionary(minWordLength)
	d.generateLetters()
	d.generateFreqs()
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

func (d *AaranagramsDictionary) generateLetters() {
	d.letters = make(map[string]int)
	for word := range d.dictionary {
		first := word[0]
		last := word[len(word)-1]
		l := string(first) + string(last)
		if _, ok := d.letters[l]; !ok {
			d.letters[l] = 1
		} else {
			d.letters[l]++
		}
	}
}

func (d *AaranagramsDictionary) generateFreqs() {
	values := make([]int, 0, len(d.letters))
	for _, v := range d.letters {
		values = append(values, v)
	}
	for i, e := range values {
		if i == 0 || e < d.minFreq {
			d.minFreq = e
		}
	}
	for i, e := range values {
		if i == 0 || e > d.maxFreq {
			d.maxFreq = e
		}
	}
}

func (d *AaranagramsDictionary) getWorth(letters string) int {
	return int(50-50*(float32(d.letters[letters]-d.minFreq)/float32(d.maxFreq-d.minFreq))) + 50
}
