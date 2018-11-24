package wfreq

import (
    "regexp"
    "strings"
    "unicode"
)

/*  input: a slice of strings that is the book
    output:
        1. a map of words and frequency of occurence of each word
        2. a slice with each element 1:1 with the lines of the text
           containing a set, per line, of words on that line
    */

func GetWordList(wb []string) (map[string]int, []map[string]struct{}) {
    f := func(c rune) bool {
        return !unicode.IsLetter(c) && !unicode.IsNumber(c)
    }
    m := make(map[string]int)  // map to hold words, counts
    var m2 []map[string]struct{}

    // hyphenated words "auburn-haired" become "auburn①haired"
    // to preserve that it is one (hyphenated) word.
    // same for single quotes within words, which are always
    // straightened within words
    var re1 = regexp.MustCompile(`(\w)\-(\w)`)
    var re2 = regexp.MustCompile(`(\w)['’](\w)`)
    for _, element := range wb {
        rtn := make(map[string]struct{},0)  // a set; empty structs take no memory
        // need to preprocess each line
        // retain [-'’] between letters
        // need this twice to handle alternates i.e. r-u-d-e
        element := re1.ReplaceAllString(element, `${1}①${2}`)
        element = re1.ReplaceAllString(element, `${1}①${2}`)
        // need this twice to handle alternates i.e. fo’c’s’le
        element = re2.ReplaceAllString(element, `${1}②${2}`)
        element = re2.ReplaceAllString(element, `${1}②${2}`)
        // all words with special characters are protected
        t := (strings.FieldsFunc(element, f))
        for _, word := range t {
            // put the special characters back in there
            s := strings.Replace(word, "①", "-", -1)
            s = strings.Replace(s, "②", "'", -1)
            // and build the map
            if _, ok := m[s]; ok {  // if it is there already, increment
               m[s] = m[s] + 1
            } else {
                m[s] = 1
            }
            rtn[s] = struct{}{}
        }
        m2 = append(m2, rtn)
    }
    return m, m2
}
