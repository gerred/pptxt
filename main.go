/*
filename:  pptxt.go
author:    Roger Frank
license:   GPL
status:    development
usage:
  ./pptxt -i --useBOM lightning-utf8.txt
  it looks for a file pptxt.dat
    -- dictionary word list
    -- jeebies
    -- scannos
  it looks for a file goodwords.txt or a user-specified filename
  it generates
    -- report file (default filename report.txt)
    -- if DEBUG: suspects.txt list of words
    -- logfile.txt report of the pptxt run, parameters used, etc.
main data structures:
  wb working buffer, one text line per slice element
  wd working dictionary, including goodwords.txt if provided
  pb paragraph buffer, one paragraph per slice element
  sw suspect words list
*/

package main

import (
    "flag"
    "fmt"
    "os"
    "pptxt/fileio"
    "pptxt/dict"
    "pptxt/leven"
    "pptxt/spellcheck"
    "pptxt/textcheck"
    "sort"
    "time"
    "path/filepath"
)

// const VERSION string = "0.90"

type Params struct {
    infile  string
    datfile string
    gwfilename string
    experimental bool
    useBOM  bool
    useCRLF bool
}

var p Params
var runlog []string  // logfile for pptxt
var report []string  // report for pptxt
var wlm []string  // suspect words returned as list by spellcheck
var wb []string  // working buffer
var wd []string  // working dictionary inc. goodwords.txt
var sw []string  // suspect words list

func Test(p Params) {
    // fmt.Println(p.experimental)
}

func doParams() Params {
    p := Params{}
    flag.StringVar(&p.infile, "i", "", "input file")
    flag.StringVar(&p.datfile, "d", "pptxt.dat", "data file")
    flag.StringVar(&p.gwfilename, "g", "goodwords.txt", "good word list")
    flag.BoolVar(&p.experimental, "x", false, "experimental (developers only)")
    flag.BoolVar(&p.useBOM, "useBOM", false, "use BOM on text output")
    flag.BoolVar(&p.useCRLF, "useCRLF", false, "CRLF line endings on output")
    flag.Parse()
    return p
}

func main() {
    // spellcheck.Debug = DEBUG
    runlog = append(runlog, fmt.Sprintf("report for pptxt\nrun started: %s",
        time.Now().Format(time.RFC850)))
    runlog = append(runlog, fmt.Sprintf("command line: %s", os.Args))

    p = doParams()  // parse command line parameters

    /*************************************************************************/
    /* working buffer (wb)                                                   */
    /* user-supplied source file UTF-8 encoded                               */
    /*************************************************************************/
    wb = fileio.ReadText(p.infile)  // working buffer, line by line

    // location of executable and user's working directory
    execut, _ := os.Executable()
    loc_exec := filepath.Dir(execut)  // i.e. /home/rfrank/go/src/pptxt
    runlog = append(runlog, fmt.Sprintf("executable is in: %s", loc_exec))
    loc_proj, _ := os.Getwd()  // i.e. /home/rfrank/projects/books/hiking-westward
    runlog = append(runlog, fmt.Sprintf("project is in: %s", loc_proj))

    /*************************************************************************/
    /* working dictionary (wd)                                               */
    /* create from words in dictionary (in pptxt.dat file)                   */
    /* and words from optional project-specific goodwords.txt file           */
    /* result are all known good words in a sorted list                      */
    /*************************************************************************/

    // default dictionary is in the pptxt.dat file
    // search in same folder as executable; if not there, search project folder
    if _, err := os.Stat(filepath.Join(loc_exec, p.datfile)); !os.IsNotExist(err) {
        // it exists
        wd = dict.ReadDict(filepath.Join(loc_exec, p.datfile))
        runlog = append(runlog, 
            fmt.Sprintf("datafile: %s",filepath.Join(loc_exec, p.datfile)))
    }
    if len(wd) == 0 {
        if _, err := os.Stat(filepath.Join(loc_proj, p.datfile)); !os.IsNotExist(err) { 
            // it exists
            wd = dict.ReadDict(filepath.Join(loc_proj, p.datfile))
            runlog = append(runlog, 
                fmt.Sprintf("datafile: %s",filepath.Join(loc_proj, p.datfile)))
        }
    }
    if len(wd) == 0 {
        runlog = append(runlog, fmt.Sprintf("no dictionary present"))
    } else {
        runlog = append(runlog, fmt.Sprintf("dictionary present: %d words", len(wd)))
    }
    wl := []string{}
    if len(p.gwfilename) > 0 { // a good word list was specified or default accepted
        if _, err := os.Stat(p.gwfilename); !os.IsNotExist(err) {  // it exists
            wl = dict.ReadWordList(filepath.Join(loc_proj,p.gwfilename))
            runlog = append(runlog, fmt.Sprintf("good word list: %d words", len(wl)))
            wd = append(wd, wl...)  // add goodwords into dictionary
        } else {  // it does not exist
            runlog = append(runlog, fmt.Sprintf("no %s found.", p.gwfilename))
        }
    }
    // need the words in a sorted list for binary search later
    if len(wl) > 0 {
        sort.Strings(wd)  // appended wordlist needs sorting
    }

    /*************************************************************************/
    /* paragraph buffer (pb)                                                 */
    /* the user source file one paragraph per line                           */
    /*************************************************************************/

    var cp string  // current (in progress) paragraph
    var pb []string  // paragraph buffer
    for _, element := range wb {
        // if this is a blank line and there is a paragraph in progress, save it
        // if not a blank line, put it into the current paragraph
        if element == "" {
            if len(cp) > 0 {
                pb = append(pb, cp) // save this paragraph
                cp = cp[:0] // empty the current paragraph buffer
            }
        } else {
            if len(cp) == 0 {
                cp += element
            } else {
                cp = cp + " " + element
            }
        }
    }
    // finished processing all lines in the file
    // flush possible non-empty current paragraph buffer    
    if len(cp) > 0 {
        pb = append(pb, cp) // save this paragraph
    }
    runlog = append(runlog, fmt.Sprintf("paragraphs: %d", len(pb)))

    /*************************************************************************/
    /* begin individual tests                                                */
    /*************************************************************************/

    // spellcheck
    // generates report in logspell.txt
    // returns list of suspect words, ok words used in text
    sw, okwords := spellcheck.Spellcheck(wb, wd, &runlog, "logspell.txt")

    // levenshtein check
    // compares all suspect words to all okwords in text
    // generates report in loglev.txt
    leven.Levencheck(wb, okwords, sw, &runlog, "loglev.txt")

    // text check
    // 
    // generates report in logtext.txt
    textcheck.Textcheck(pb, wb, &runlog, "logtext.txt")    

    /*************************************************************************/
    /* all tests complete. save results to specified report file and logfile */
    /*************************************************************************/

    fileio.SaveText(runlog, "logpptxt.txt", p.useBOM, p.useCRLF)

    // remaining words in sw are suspects. conditionally generate a report
    var s []string
    if p.experimental {
        for _, word := range sw {
           s = append(s, fmt.Sprintf("%s", word))
        }
        fileio.SaveText(s, "logsuspects.txt", p.useBOM, p.useCRLF)
    }

    Test(p)
}
