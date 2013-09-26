// Porter Stemmer Algorithm in Go
// Developed by Suhail Patel <me@suhailpatel.com>
//
// Copyright (C) 2013 Suhail Patel
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this
// software and associated documentation files (the "Software"), to deal in the Software
// without restriction, including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons
// to whom the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or
// substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
// PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR
// ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// This is an implementation of the Porter Stemmer algorithm by Martin Porter as
// described in Chapter 6 of the report:
//     C.J. van Rijsbergen, S.E. Robertson and M.F. Porter, 1980. New models in
//     probabilistic information retrieval. London: British Library. (British
//     Library Research and Development Report, no. 5587).
//
// The implementation is improved slightly from the original paper as recommended
// by Martin Porter in the section 'Points of difference from the published algorithm'
// on the link below
//
// See http://www.tartarus.org/~martin/PorterStemmer for more information about
// the algorithm
package stemmer

import (
	"strings"
)

// Stem takes in a word and runs through the various steps of the original Porter Stemmer
// algorithm. This method only works on words in English and expects all words to be sanitized
// (trimmed etc.). Words of length 1 or 2 are ignored (as recommended by the points of
// difference of the algorithm). Words are converted to lower case and will be returned in
// lower case
func Stem(word string) string {
	if len(word) == 1 || len(word) == 2 {
		return word
	}

	stemmed := strings.TrimSpace(word)
	stemmed = strings.ToLower(stemmed)
	stemmed = step1a(stemmed)
	stemmed = step1b(stemmed)
	stemmed = step1c(stemmed)
	stemmed = step2(stemmed)
	stemmed = step3(stemmed)
	stemmed = step4(stemmed)
	stemmed = step5(stemmed)

	return stemmed
}

// Step 1A focuses on getting rid of plurals
func step1a(word string) string {
	matched := false

	word, matched = checkReplace(word, "sses", "ss", matched, nil)
	word, matched = checkReplace(word, "ies", "i", matched, nil)
	word, matched = checkReplace(word, "ss", "ss", matched, nil)
	word, matched = checkReplace(word, "s", "", matched, nil)

	return word
}

// Step 1B also focuses on getting rid of plurals
func step1b(word string) string {
	matched := false

	word, matched = checkReplace(word, "eed", "ee", matched, func(stem string) bool {
		return m(stem) > 0
	})

	vowel := func(stem string) bool {
		return hasVowel(stem)
	}

	if matched == true {
		return word
	}

	prevWord := word
	word, matched = checkReplace(word, "ed", "", matched, vowel)
	word, matched = checkReplace(word, "ing", "", matched, vowel)

	if matched == true && word != prevWord {
		matched = false
		word, matched = checkReplace(word, "at", "ate", matched, nil)
		word, matched = checkReplace(word, "bl", "ble", matched, nil)
		word, matched = checkReplace(word, "iz", "ize", matched, nil)

		i := len(word)
		if matched == false && i >= 2 && hasDoubleConsonantSuffix(word[0:i]) &&
			word[i-2] != 'l' && word[i-2] != 's' && word[i-2] != 'z' {
			word, matched = word[0:i-1], true
		}

		word, matched = checkReplace(word, "", "e", matched, func(stem string) bool {
			return m(stem) == 1 && cvc(stem)
		})
	}

	return word
}

// Focus on past principles
func step1c(word string) string {
	word, _ = checkReplace(word, "y", "i", false, func(stem string) bool {
		return hasVowel(stem)
	})

	return word
}

// Step 2 is just going through the rules, switch on the penultimate letter
// for a speed boost in comparison for which rules match/execute
func step2(word string) string {
	matched := false

	mMoreZero := func(stem string) bool {
		return m(stem) > 0
	}

	switch word[len(word)-2] {
	case 'a':
		word, matched = checkReplace(word, "ational", "ate", matched, mMoreZero)
		word, matched = checkReplace(word, "tional", "tion", matched, mMoreZero)
		break

	case 'c':
		word, matched = checkReplace(word, "enci", "ence", matched, mMoreZero)
		word, matched = checkReplace(word, "anci", "ance", matched, mMoreZero)
		break

	case 'e':
		word, matched = checkReplace(word, "izer", "ize", matched, mMoreZero)
		break

	case 'g':
		word, matched = checkReplace(word, "logi", "log", matched, mMoreZero)
		break

	case 'l':
		word, matched = checkReplace(word, "bli", "ble", matched, mMoreZero)
		word, matched = checkReplace(word, "alli", "al", matched, mMoreZero)
		word, matched = checkReplace(word, "entli", "ent", matched, mMoreZero)
		word, matched = checkReplace(word, "eli", "e", matched, mMoreZero)
		word, matched = checkReplace(word, "ousli", "ous", matched, mMoreZero)
		break

	case 'o':
		word, matched = checkReplace(word, "ization", "ize", matched, mMoreZero)
		word, matched = checkReplace(word, "ation", "ate", matched, mMoreZero)
		word, matched = checkReplace(word, "ator", "ate", matched, mMoreZero)
		break

	case 's':
		word, matched = checkReplace(word, "alism", "al", matched, mMoreZero)
		word, matched = checkReplace(word, "iveness", "ive", matched, mMoreZero)
		word, matched = checkReplace(word, "fulness", "ful", matched, mMoreZero)
		word, matched = checkReplace(word, "ousness", "ous", matched, mMoreZero)
		break

	case 't':
		word, matched = checkReplace(word, "aliti", "al", matched, mMoreZero)
		word, matched = checkReplace(word, "iviti", "ive", matched, mMoreZero)
		word, matched = checkReplace(word, "biliti", "ble", matched, mMoreZero)
		break

	default:
		break
	}

	return word
}

// Stemming words as part of Step 3
func step3(word string) string {
	mMoreZero := func(stem string) bool {
		return m(stem) > 0
	}

	matched := false

	word, matched = checkReplace(word, "icate", "ic", matched, mMoreZero)
	word, matched = checkReplace(word, "ative", "", matched, mMoreZero)
	word, matched = checkReplace(word, "alize", "al", matched, mMoreZero)
	word, matched = checkReplace(word, "iciti", "ic", matched, mMoreZero)
	word, matched = checkReplace(word, "ical", "ic", matched, mMoreZero)
	word, matched = checkReplace(word, "ful", "", matched, mMoreZero)
	word, matched = checkReplace(word, "ness", "", matched, mMoreZero)

	return word
}

// More stemming as part of Step 4
func step4(word string) string {
	mMoreOne := func(stem string) bool {
		return m(stem) > 1
	}

	matched := false

	word, matched = checkReplace(word, "al", "", matched, mMoreOne)
	word, matched = checkReplace(word, "ance", "", matched, mMoreOne)
	word, matched = checkReplace(word, "ence", "", matched, mMoreOne)
	word, matched = checkReplace(word, "er", "", matched, mMoreOne)
	word, matched = checkReplace(word, "ic", "", matched, mMoreOne)
	word, matched = checkReplace(word, "able", "", matched, mMoreOne)
	word, matched = checkReplace(word, "ible", "", matched, mMoreOne)
	word, matched = checkReplace(word, "ant", "", matched, mMoreOne)
	word, matched = checkReplace(word, "ement", "", matched, mMoreOne)
	word, matched = checkReplace(word, "ment", "", matched, mMoreOne)
	word, matched = checkReplace(word, "ent", "", matched, mMoreOne)

	word, matched = checkReplace(word, "ion", "", matched, func(stem string) bool {
		return m(stem) > 1 && (stem[len(stem)-1] == 's' || stem[len(stem)-1] == 't')
	})

	word, matched = checkReplace(word, "ou", "", matched, mMoreOne)
	word, matched = checkReplace(word, "ism", "", matched, mMoreOne)
	word, matched = checkReplace(word, "ate", "", matched, mMoreOne)
	word, matched = checkReplace(word, "iti", "", matched, mMoreOne)
	word, matched = checkReplace(word, "ous", "", matched, mMoreOne)
	word, matched = checkReplace(word, "ive", "", matched, mMoreOne)
	word, matched = checkReplace(word, "ize", "", matched, mMoreOne)

	return word
}

// Step5 focuses on clean up (the paper splits it up to A and B but i've
// combined it )
func step5(word string) string {
	matched := false
	word, matched = checkReplace(word, "e", "", matched, func(stem string) bool {
		return m(stem) > 1 || m(stem) == 1 && !cvc(stem)
	})

	matched = false
	word, _ = checkReplace(word, "l", "", matched, func(stem string) bool {
		return m(stem) > 1 && stem[len(stem)-1] == 'l' && hasDoubleConsonantSuffix(word)
	})

	return word
}

// Determines whether a character at the position specified is a consonant
// as defined by the paper. A \consonant\  in a word is a letter other than
// A, E, I, O or U, and other than Y preceded by a consonant.
func consonant(word string, position int) bool {
	i := position

	if word[i] == 'a' || word[i] == 'e' || word[i] == 'i' || word[i] == 'o' || word[i] == 'u' {
		return false
	}

	if word[i] == 'y' && i > 0 {
		return vowel(word, position-1)
	}

	return true
}

// Determines whether a character at the position specified is a vowel which
// is a inverse of consonant
func vowel(word string, position int) bool {
	return !consonant(word, position)
}

// Loops through the word specified and determines whether it has a vowel in
// the word as defined in the paper
func hasVowel(word string) bool {
	for i := 0; i < len(word); i++ {
		if vowel(word, i) {
			return true
		}
	}

	return false
}

// Determines whether the word specified has a double same-character consonant
// at the end (consonant as defined by the paper)
func hasDoubleConsonantSuffix(word string) bool {
	i := len(word) - 1
	if word[i] == word[i-1] && consonant(word, i) {
		return true
	}

	return false
}

// A word has the cvc form if the last 3 characters follows the
// consonant - vowel - consonant structure and the second consonant is
// not w, x or y
func cvc(word string) bool {
	length := len(word)

	if length < 3 {
		return false
	}

	if consonant(word, length-1) && vowel(word, length-2) && consonant(word, length-3) {
		if word[length-1] == 'w' || word[length-1] == 'x' || word[length-1] == 'y' {
			return false
		}
		return true
	}
	return false
}

// Determines whether the word has suffix string as it's suffix
func hasSuffix(word string, suffix string) bool {
	if len(suffix) == 0 {
		return true
	}

	if len(word) < len(suffix) {
		return false
	}

	return strings.HasSuffix(word, suffix)
}

// Condition function type where we can specify and stem and have a
// condition function to determine whether the stem matches any criteria
// it needs to match
type stemCondition func(stem string) bool

// This method determines whether a word can be matched by removing it's old suffix
// and tacking on the replacement. The matched bool parameter is used to determine whether
// a replacement has already taken place (in which case don't do anything). Note, even if
// the stem condition fails, it does return true for matched because the matched
// is only used to determine whether the suffix matched, not if a replacement was made
func checkReplace(word string, suffix string, replace string, matched bool, condition stemCondition) (string, bool) {
	if matched || len(word) < len(suffix) {
		return word, matched
	}

	if len(suffix) > 0 && !hasSuffix(word, suffix) {
		return word, matched
	}

	if condition != nil && !condition(word[0:len(word)-len(suffix)]) {
		return word, true
	}

	return word[:len(word)-len(suffix)] + replace, true
}

// m() measures the number of consonant sequences between the start and end of the
// word provided. c denotes a consonant sequence and v a vowel sequence, and <..>
// indicates arbitrary presence
//
//      <c><v>       gives 0
//      <c>vc<v>     gives 1
//      <c>vcvc<v>   gives 2
//      <c>vcvcvc<v> gives 3
//      ....
func m(word string) int {
	length := len(word)
	charsSeen := 0
	mCount := 0

	// Loop through all the initial consonants first
	for {
		if charsSeen >= length {
			return mCount
		}
		if vowel(word, charsSeen) {
			break
		}
		charsSeen++
	}
	charsSeen++

	// Now look for the VC{m} pairs
	for {
		for {
			if charsSeen >= length {
				return mCount
			}
			if consonant(word, charsSeen) {
				break
			}
			charsSeen++
		}
		charsSeen++
		mCount++

		for {
			if charsSeen >= length {
				return mCount
			}
			if vowel(word, charsSeen) {
				break
			}
			charsSeen++
		}
		charsSeen++
	}

	return mCount
}
