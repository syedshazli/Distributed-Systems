package cs4513_go_impl

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

/*
Find the top K most common words in a text document.
What is a word?

	A word here only consists of alphanumeric characters, e.g., catch21
	All punctuations and other characters should be removed, e.g. "don't" becomes "dont" or "end." becomes "end"; done before the charThreshold
	A word has to satifies the charThreshold, e.g., if charThreshold = 5  "apple" is a word, but neither "new" or "york" are words

Matching condition

	Matching is case insensitive

Parameters:
- path: file path
- numWords: number of words to return (i.e. k)
- charThreshold: threshold for whether a token qualifies as a word
You should use `checkError` to handle potential errors.
*/
func TopWords(path string, numWords int, charThreshold int) []WordCount {

	// HINT: You may find the `strings.Fields` and `strings.ToLower` functions helpful.
	// HINT: the regex "[^0-9a-zA-Z]+" can be used to spot any non-alphanumeric characters.

	if numWords < 1 {
		return []WordCount{}
	}
	// read in file into string
	// the whole file is now in the 'content' var
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	contentToString := string(content)

	// the content of the file is now in lowercase in this string
	contentLowerCase := strings.ToLower(contentToString)
	rawTokens := strings.Fields(contentLowerCase) // each element in raawtokens is a whitespaced split word

	var filtered []string

	punctuationRegex := regexp.MustCompile(`[^a-zA-Z0-9]+`)

	for _, token := range rawTokens {

    	cleanWord := punctuationRegex.ReplaceAllString(token, "") // remove punctuation, replace w/nothing
    
		// must be greater than or equal to threshold to be counted
    	if len(cleanWord) >= charThreshold {
        	filtered = append(filtered, cleanWord)
    	}
}

	/*
		1. remove all non-alphanumeric chars from the string ✅
		2. Find all words in the string ✅
		3. See how many times a unique word shows up, counting the word only if len is >= charThreshold ✅
		4. Make a list of all unique words and their counts, sort by count ✅
		5. Return/iterate up to k ✅
	*/

	

	// words is now an array of words that is of len greater than or equal to the threshold
	words := filtered
	fmt.Println(words)

	wordMap := make(map[string]int)

	// iterate thru each word in the string. if the word is in the map, increment the value by 1. if not, init with 1
	for _, word := range words {
		value, ok := wordMap[word]
		if ok {
			wordMap[word] = value + 1
		} else {
			wordMap[word] = 1
		}
	}

	arrayOfWordCounts := []WordCount{}

	for wordStr, mapCount := range wordMap {
		newWord := WordCount{Word: wordStr, Count: mapCount}
		arrayOfWordCounts = append(arrayOfWordCounts, newWord)
	}

	// this now has an array of structs containing words and their counts, sorted by count in descending order
	sortWordCounts(arrayOfWordCounts)

	// now, remove any element not in the top 'numWords' elements
	if len(arrayOfWordCounts) < numWords {
		numWords = len(arrayOfWordCounts)
	}
	arrayOfWordCounts = arrayOfWordCounts[:numWords]

	return arrayOfWordCounts

}

/*
Do NOT modify this struct!
A struct that represents how many times a word is observed in a document
*/
type WordCount struct {
	Word  string
	Count int
}

/*
Do NOT modify this function!
*/
func (wc WordCount) String() string {
	return fmt.Sprintf("%v: %v", wc.Word, wc.Count)
}

/*
Do NOT modify this function!
Helper function to sort a list of word counts in place.
This sorts by the count in decreasing order, breaking ties using the word.
*/
func sortWordCounts(wordCounts []WordCount) {
	sort.Slice(wordCounts, func(i, j int) bool {
		wc1 := wordCounts[i]
		wc2 := wordCounts[j]
		if wc1.Count == wc2.Count {
			return wc1.Word < wc2.Word
		}
		return wc1.Count > wc2.Count
	})
}
