package textcheck

import (
    "fmt"
    "pptxt/fileio"
    "time"
    "strings"
    "strconv"
    "sort"
)

var s []string  // to build the log specific to this test
var rs []string  // to append to the overall runlog for all tests

func report(r string) {
    s = append(s, r)
}

func asteriskCheck(wb []string) {
    report("asterisk check")
    count := 0
    for n, line := range wb {
        if strings.Contains(line, "*") {
            report(fmt.Sprintf("  %d: %s", n, line))
            count += 1
        }
    }
    if count == 0 {
        report("  no unexpected asterisks found in text.")
    }
}

// do not report adjacent spaces that start or end a line
func adjacentSpaces(wb []string) {
    report("adjacent spaces check")
    count := 0
    for n, line := range wb {
        if strings.Contains(strings.TrimSpace(line), "  ") {
            report(fmt.Sprintf("  %d: %s", n, line))
            count += 1
        }
    }
    if count == 0 {
        report("  no adjacent spaces found in text.")
    }
}

// 
func trailingSpaces(wb []string) {
    report("trailing spaces check")
    count := 0
    for n, line := range wb {
        if strings.TrimSuffix(line, " ") != line {
            report(fmt.Sprintf("  %d: %s", n, line))
            count += 1
        }
    }
    if count == 0 {
        report("  no trailing spaces found in text.")
    }
}

type kv struct {
    Key   rune
    Value int
}

var m = map[rune]int{}  // a map for letter frequency counts

// report infrequently-occuring characters (runes)
// threshold set to fewer than 10 occurences
// do not report numbers
func letterChecks(wb []string) {
    report("character checks")
    count := 0
    for _, line := range wb {
        for _, char := range(line) {  // this gets runes
            m[char] += 1 
        }
    }
    var ss []kv // slice of Key, Value pairs
    for k, v := range m {  // load it up
        ss = append(ss, kv{k, v})
    }
    sort.Slice(ss, func(i, j int) bool {  // sort it based on Value
        return ss[i].Value > ss[j].Value
    })
    for _, kv := range ss {
        reportme := false
        if kv.Value < 10 && (kv.Key < '0' || kv.Key > '9') {
            reportme = true
        }
        if !strings.ContainsRune(",:;—?!-_0123456789“‘’”. abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", kv.Key) {
            reportme = true
        }
        if reportme {
            reportcount := 0
            report(fmt.Sprintf("  %s", strconv.QuoteRune(kv.Key)))
            count += 1
            for n, line := range wb {
                if strings.ContainsRune(line, kv.Key) {
                    if reportcount < 5 {
                        report(fmt.Sprintf("    %d: %s", n, line))
                    }
                    if reportcount == 5 {
                        report(fmt.Sprintf("    ...more"))
                    }
                    reportcount++
                }
            }
        }
    }
    if count == 0 {
        report("  no character checks reported.")
    }       
}

// special situations only report if they find something
func spacingCheck(wb []string) {
    count := 0
    s := ""
    report("spacing check")
    
    consec := 0  // consecutive blank lines
    for _,line := range wb {
        if line == "" {
            consec++
        } else { // a non-blank line
            if consec >= 4 {  // start of a new chapter
                fmt.Println(s)
                s = "4"
            } else {
                if consec > 0 {
                    s = fmt.Sprintf("%s%d", s, consec)
                }
            }
            consec = 0
        }
    }

    if count == 0 {
        report("  no spacing errors reported.")
    }       
}

// special situations only report if they find something
func specialSituations(wb []string) {
    count := 0
    report("special situations checks")
    
    if m['\''] > 0 && ( m['‘'] > 0 || m['’'] > 0 ) {
        report("  both straight and curly single quotes found in text")
        count++
    }
    if m['"'] > 0 && ( m['“'] > 0 || m['”'] > 0 ) {
        report("  both straight and curly double quotes found in text")
        count++
    }

    if count == 0 {
        report("  no special situations checks reported.")
    }       
}

// text checks
// a series of tests either on the working buffer (line at a time)
// or the paragraph buffer (paragraph at a time)
func Textcheck(pb []string, wb []string, runlog *[]string, fname string) {
    rs = append(rs, "Text checks")
    s = append(s, fmt.Sprintf("text check report \nstarted: %s\n------------------------------",
        time.Now().Format(time.RFC850)))

    rs = append(rs, "  added to runlog by Textcheck")

    asteriskCheck(wb)
    adjacentSpaces(wb)
    trailingSpaces(wb)
    letterChecks(wb)
    spacingCheck(wb)
    specialSituations(wb)

    // generate logtext.txt from s
    fileio.SaveText(s, fname, true, true)

    // append to pptxt.log
    *runlog = append(*runlog, rs...)
}