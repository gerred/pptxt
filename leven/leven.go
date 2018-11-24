package leven

import (
	"fmt"
	"github.com/asylumcs/pptxt/fileio"
	"github.com/asylumcs/pptxt/wfreq"
	"strings"
	"time"
)

func Levenshtein(str1, str2 []rune) int {
	s1len := len(str1)
	s2len := len(str2)
	column := make([]int, len(str1)+1)

	for y := 1; y <= s1len; y++ {
		column[y] = y
	}
	for x := 1; x <= s2len; x++ {
		column[0] = x
		lastkey := x - 1
		for y := 1; y <= s1len; y++ {
			oldkey := column[y]
			var incr int
			if str1[y-1] != str2[x-1] {
				incr = 1
			}

			column[y] = minimum(column[y]+1, column[y-1]+1, lastkey+incr)
			lastkey = oldkey
		}
	}
	return column[s1len]
}

func minimum(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else {
		if b < c {
			return b
		}
	}
	return c
}

// iterate over every suspect word at least six letters long
// case insensitive
// looking for a good word in the text that is "near"
func Levencheck(wb []string, okwords []string, suspects []string, runlog *[]string, fname string) {
	var s []string
	var rs []string
    rs = append(rs, "Levenshtein checks")

	s = append(s, fmt.Sprintf("distance check report \nstarted: %s\n------------------------------",
		time.Now().Format(time.RFC850)))
    count := 0

    var wb2 []map[string]struct{}
    _, wb2 = wfreq.GetWordList(wb) // wordlist with word, frequency of word in map, wb2

	for _, suspect := range suspects {
		for _, okword := range okwords {
			if len(suspect) < 6 {
				continue
			}
			if strings.ToLower(suspect) == strings.ToLower(okword) {
				continue
			}
			dist := Levenshtein([]rune(suspect),[]rune(okword))
			if dist < 2 {
				// get counts
				suspectwordcount := 0
				okwordcount := 0
		        for n, _ := range(wb) {
		            wordsonthisline := wb2[n] // a set of words on this line
		            if _, ok := wordsonthisline[suspect]; ok {
		            	suspectwordcount += 1
		            }
		            if _, ok := wordsonthisline[okword]; ok {
		            	okwordcount += 1
		            }
		        }
				s = append(s, fmt.Sprintf("%s(%d):%s(%d)", suspect, suspectwordcount, okword, okwordcount))

				// show one line in context
		        count := 0
		        for n, line := range(wb) {
		            wordsonthisline := wb2[n] // a set of words on this line
		            if _, ok := wordsonthisline[suspect]; ok {
		            	if count == 0 {
		            		s = append(s, fmt.Sprintf("  %6d: %s", n, line))	
		            	}
		                count += 1
		            }
		        }
		        count = 0
		        for n, line := range(wb) {
		            wordsonthisline := wb2[n] // a set of words on this line
		            if _, ok := wordsonthisline[okword]; ok {
		            	if count == 0 {	            	
		                	s = append(s, fmt.Sprintf("  %6d: %s", n, line))
		                }
		                count += 1
		            }
		        }
			}
		}
	}

    rs = append(rs, fmt.Sprintf("  suspect words by distance check: %d", count))

    // generate loglev.txt from s
    fileio.SaveText(s, fname, true, true)

    // append to pptxt.log
    *runlog = append(*runlog, rs...)
}