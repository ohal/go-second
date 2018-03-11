package main

import (
	"fmt"
	//"log"
	//"regexp"
	"sort"
	"strings"
	//"unicode"
	//"github.com/mfonda/simhash"
)

// ReviewPost data
type ReviewPost struct {
	ID   int    `json:"id"`
	Post string `json:"post"`
}

// StrShingleFreq data
type StrShingleFreq struct {
	Freq    int    `json:"freq"`
	Shingle string `json:"shingle"`
}

func sortPhrasesByFreq(list []string) []StrShingleFreq {
	sMap := shingleCount(list)
	freq := convertShingles(sMap)
	sort.Sort(byFreq(freq))

	return freq
}

/*
func sortPhrasesByPhrase(list []string) []StrShingleFreq {
    sMap := shingleCount(list)
    freq := convertShingles(sMap)
    sort.Sort(byPhrase(freq))

    return freq
}
*/

func getShinglesList(post ReviewPost) []string {

	p := strings.ToLower(post.Post)

	// TODO: split to words with/without punctuation
	/*
	   f := func(c rune) bool {
	       return !unicode.IsLetter(c) && !unicode.IsNumber(c) && !unicode.IsPunct(c)
	   }

	   l := strings.FieldsFunc(p, f)
	*/

	removePunctuation := func(r rune) rune {
		if strings.ContainsRune(".,:;!()-", r) {
			return -1
		}
		return r
	}

	p = strings.Map(removePunctuation, p)
	words := strings.Fields(p)

	ph := make([]string, 0)
	for i := 0; i < len(words)-shingleSize+1; i++ {
		sw := words[i : i+shingleSize]
		shingle := strings.Join(sw, " ")
		ph = append(ph, shingle)
	}

	// give all possible shingles from word list
	/*
	   re := regexp.MustCompile("\\w+")
	   l := re.FindAllString(p, -1)

	   b := make([][]byte, 0)
	   for _, w := range l {
	       b = append(b, []byte(w))
	   }

	   // get shingles from word list
	   s := simhash.Shingle(3, b)

	   ph := make([]string, 0)
	   for _, bs := range s {
	       ph = append(ph, string(bs))
	   }
	*/

	return ph
}

/*
func getShinglesListCh(post ReviewPost) chan string {

    p := strings.ToLower(post.Post)

    removePunctuation := func(r rune) rune {
        if strings.ContainsRune(".,:;!()-", r) {
            return -1
        }
        return r
    }

    p = strings.Map(removePunctuation, p)
    words := strings.Fields(p)

    out := make(chan string)
    go func() {
        for i := 0; i < len(words)-2; i++ {
            sw := words[i : i+3]
            shingle := strings.Join(sw, " ")
            out <- shingle
        }
        close(out)
    }()

    return out
}
*/

func shingleCount(shingles []string) map[string]int {
	shingleMap := make(map[string]int)
	for _, shingle := range shingles {
		if _, ok := shingleMap[shingle]; ok {
			shingleMap[shingle]++
		} else {
			shingleMap[shingle] = 1
		}
	}
	return shingleMap
}

// StrShingleFreq will be displayed in this format
func (s StrShingleFreq) String() string {
	return fmt.Sprintf("%3d   %s", s.Freq, s.Shingle)
}

func convertShingles(mapPhrases map[string]int) []StrShingleFreq {
	// convert the map to a slice of structures for sorting
	// create pointer to an instance of shingle + frequency
	p := new(StrShingleFreq)
	slice := make([]StrShingleFreq, len(mapPhrases))
	ix := 0
	for key, value := range mapPhrases {
		p.Shingle, p.Freq = key, value
		slice[ix] = *p
		ix++
	}

	return slice
}

// byFreq implements sort.Interface for []strShingleFreq based on the freq field
type byFreq []StrShingleFreq

func (a byFreq) Len() int           { return len(a) }
func (a byFreq) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byFreq) Less(i, j int) bool { return a[i].Freq < a[j].Freq }

/*
// byPhrase implements sort.Interface for []strShingleFreq based on the word field
type byPhrase []StrShingleFreq

func (a byPhrase) Len() int           { return len(a) }
func (a byPhrase) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPhrase) Less(i, j int) bool { return a[i].Shingle < a[j].Shingle }
*/
