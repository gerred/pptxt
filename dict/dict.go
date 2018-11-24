package dict

import (
    "os"
    "log"
    "bufio"
    "strings"
)

var BOM = string([]byte{239, 187, 191}) // UTF-8 specific

// dictionary word list in in pptxt.dat bracketed by
// *** BEGIN DICT *** and *** END DICT ***
func ReadDict(infile string) []string {
    file, err := os.Open(infile)
    if err != nil { log.Fatal(err) }
    defer file.Close()
    wd := []string{}
    scanner := bufio.NewScanner(file)
    keep := false
    for scanner.Scan() {
        if scanner.Text() == "*** BEGIN DICT ***" {
            keep = true
            continue    
        }
        if scanner.Text() == "*** END DICT ***" {
            keep = false
            continue
        }
        if keep {
            wd = append(wd, scanner.Text())
        }
    }
    if err := scanner.Err(); err != nil { log.Fatal(err) }

    // remove BOM if present
    wd[0] = strings.TrimPrefix(wd[0], BOM)
    return wd
}

func ReadWordList(infile string) []string {
    wd := []string{}
    file, err := os.Open(infile)  // try to open wordlist
    if err != nil {
        return wd  // early exit if it isn't present
    }
    defer file.Close()  // here if it opened
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        wd = append(wd, scanner.Text())
    }
    if err := scanner.Err(); err != nil { log.Fatal(err) }
    // remove BOM if present
    if len(wd) > 0 {
        wd[0] = strings.TrimPrefix(wd[0], BOM)  // on first word if there is one
    }
    return wd
}
