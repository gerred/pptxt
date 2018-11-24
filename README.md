# pptxt
post-processor's text validation tool to Project Gutenberg standards

replaces all individual Python3 versions for text files in the
Post Processor's Workbench.

compiled, written in the Go language.
platform independent binaries

using from a working book folder:
rfrank@carbon:~/projects/books/hiking-westward
$ (cd ~/go/src/pptxt && go build) && ~/go/src/pptxt/pptxt -i westward-utf8.txt
