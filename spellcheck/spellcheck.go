package spellcheck

import (
    "github.com/asylumcs/pptxt/fileio"
    "github.com/asylumcs/pptxt/wfreq"
    "strings"
	"strconv"
    "fmt"
    "sort"
    "time"
)

func lookup(wd []string, word string) bool {
    ip := sort.SearchStrings(wd, word) // where it would insert
    return (ip != len(wd) && wd[ip] == word) // true if we found it
}

// spellcheck returns list of suspect words, list of ok words in text
func Spellcheck(wb []string, wd []string, runlog *[]string, fname string) ([]string, []string) {
    var rs []string  // for logfile.txt
    rs = append(rs, "spellcheck")

    // for speed, build wb2 to hold a map of words on each line.
    var wb2 []map[string]struct{}
    okwordlist := make(map[string]int)  // cumulative words OK by successive tests
    var willdelete []string  // words to be deleted from wordlist
    wlm, wb2 := wfreq.GetWordList(wb) // wordlist with word, frequency of word in map, wb2
    rs = append(rs, fmt.Sprintf("  unique words in text: %d words", len(wlm)))

    for word, count := range wlm {
        ip := sort.SearchStrings(wd, word) // where it would insert
        if ip != len(wd) && wd[ip] == word { // true if we found it
            // ok by wordlist
            okwordlist[word] =  count // remember as good word
            willdelete = append(willdelete,word)
        }
    }

    rs = append(rs, fmt.Sprintf("  approved by dictionary: %d words", len(willdelete)))

    // fmt.Printf("%+v\n", willdelete)
    // os.Exit(1)

    // delete words that have been OKd by being in the dictionary
    for _, word := range(willdelete) {
        delete(wlm, word)
    }
    willdelete = nil  // clear the list of words to delete

    // fmt.Println(len(wlm))
    // fmt.Println(len(okwordlist))

    // typically at this point, I have taken the 8995 unique words in the book
    // and sorted it into 7691 words found in the dictionary and 1376 words unresolved

    // fmt.Printf("%+v\n", wlm)
    // try to approve words that are capitalized by testing them lower case

    // check words ok by depossessive
    lcwordlist := make(map[string]int)  // words OK by their lowercase form being in the dictionary

    for word, count := range(wlm) {
        lcword := strings.ToLower(word)
        ip := sort.SearchStrings(wd, lcword) // where it would insert
        if ip != len(wd) && wd[ip] == lcword { // true if we found it
            // ok by lowercase
            lcwordlist[word] =  count // remember (uppercase versions) as good word
            okwordlist[word] =  count // remember as good word
            willdelete = append(willdelete,word)
        }
    }

    rs = append(rs, fmt.Sprintf("  approved by lowercase form: %d words", len(willdelete)))

    // delete words that have been OKd by their lowercase form being in the dictionary
    for _, word := range(willdelete) {
        delete(wlm, word)
    }
    willdelete = nil  // clear the list of words to delete

    // fmt.Printf("%+v\n", wlm)
    // fmt.Println(len(wlm))
    // fmt.Println(len(lcwordlist))

    // typically at this point the 1376 unresolved words are now 638 that were approved b/c their
    // lowercase form is in the dictionary and 738 words still unresolved

    // some of these are hyphenated. Break those words on hyphens and see if all the individual parts
    // are valid words. If so, approve the hyphenated version

    hywordlist := make(map[string]int)  // hyphenated words OK by all parts being words

    for word, count := range(wlm) {
        t := strings.Split(word, "-")
        if len(t) > 1 {
            // we have a hyphenated word
            allgood := true
            for _, hpart := range(t) {
                if !lookup(wd, hpart) {
                    allgood = false
                }
            }
            if allgood { // all parts of the hyhenated word are words
                hywordlist[word] = count
                okwordlist[word] =  count // remember as good word
                willdelete = append(willdelete, word)
            }
        }
    }

    rs = append(rs, fmt.Sprintf("  approved by dehyphenation: %d words", len(willdelete)))

    // delete words that have been OKd by dehyphenation
    for _, word := range(willdelete) {
        delete(wlm, word)
    }
    willdelete = nil  // clear the list of words to delete

    // fmt.Println(len(wlm))
    // fmt.Println(len(hywordlist))

    // of the 738 unresolved words before dehyphenation checks, now an additional
    // 235 have been approved and 503 remain unresolved

    // some "words" are entirely numerals. approve those
    for word, _ := range(wlm) {
        if _, err := strconv.Atoi(word); err == nil {
            okwordlist[word] =  1 // remember as good word
            willdelete = append(willdelete, word)
       }
    }

    rs = append(rs, fmt.Sprintf("  approved pure numerics: %d words", len(willdelete)))
    // delete words that are entirely numeric
    for _, word := range(willdelete) {
        delete(wlm, word)
    }
    willdelete = nil  // clear the list of words to delete

    // the 503 unresolved words are now 381 with the removal of the all-numeral words

    frwordlist := make(map[string]int)  // words ok by frequency occuring 4 or more times

    // some words occur many times. Accept them by frequency if they appear four or more times
    // spelled the same way
    for word, count := range(wlm) {
        if count >= 4 {
            frwordlist[word] = count
            okwordlist[word] =  count // remember as good word
            willdelete = append(willdelete, word)
       }
    }

    rs = append(rs, fmt.Sprintf("  approved by frequency: %d words", len(willdelete)))
    // delete words approved by frequency
    for _, word := range(willdelete) {
        delete(wlm, word)
    }
    willdelete = nil  // clear the list of words to delete

    // show each word in context
    var s []string
    var sw []string
    s = append(s, fmt.Sprintf("spellcheck report \nstarted: %s\n------------------------------",
        time.Now().Format(time.RFC850)))
    for word, _ := range(wlm) {
        sw = append(sw, word)  // simple slice of only the word
        s = append(s, fmt.Sprintf("%s", word))  // word we will show in context
        // show word in text
        for n, line := range(wb) {
            wordsonthisline := wb2[n] // a set of words on this line
            if _, ok := wordsonthisline[word]; ok {
                s = append(s, fmt.Sprintf("  %6d: %s", n, line))
            }
            // fmt.Printf("%+v\n", wordsonthisline)
            // s = append(s, fmt.Sprintf("  %d:  %s", n, line))
        }
        s = append(s, "")
    }

    rs = append(rs, fmt.Sprintf("  good words in text: %d words", len(okwordlist)))
    rs = append(rs, fmt.Sprintf("  suspect words in text: %d words", len(sw)))

    // generate logspell.txt from s
    fileio.SaveText(s, fname, true, true)

    // append to pptxt.log
    *runlog = append(*runlog, rs...)

    var ok []string
    for word, _ := range(okwordlist) {
        ok = append(ok, word)
    }

    // return sw: list of suspect words and ok: list of good words in text
    return sw, ok
}
