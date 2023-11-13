package utility

import "math/rand"

var letterFrequencies map[rune]int = make(map[rune]int)

var bag []rune

var bagSize int

func initLetterFrequencies() {
	letterFrequencies['E'] = 11
	letterFrequencies['A'] = 9
	letterFrequencies['R'] = 8
	letterFrequencies['I'] = 8
	letterFrequencies['O'] = 7
	letterFrequencies['T'] = 7
	letterFrequencies['N'] = 7
	letterFrequencies['S'] = 6
	letterFrequencies['L'] = 6
	letterFrequencies['C'] = 5
	letterFrequencies['U'] = 4
	letterFrequencies['D'] = 3
	letterFrequencies['P'] = 3
	letterFrequencies['M'] = 3
	letterFrequencies['H'] = 3
	letterFrequencies['G'] = 2
	letterFrequencies['B'] = 2
	letterFrequencies['F'] = 2
	letterFrequencies['Y'] = 2
	letterFrequencies['W'] = 1
	letterFrequencies['K'] = 1
	letterFrequencies['V'] = 1
	letterFrequencies['X'] = 1
	letterFrequencies['Z'] = 1
	letterFrequencies['J'] = 1
	letterFrequencies['Q'] = 1

	bagSize = 0
	for _, freq := range letterFrequencies {
		bagSize = bagSize + freq
	}

	bag = make([]rune, bagSize)
	c := 0
	for letter, freq := range letterFrequencies {
		for i := 0; i < freq; i++ {
			bag[c] = letter
			c++
		}
	}
}

func GetLetterWeighted() rune {
	return bag[rand.Intn(bagSize)]
}
