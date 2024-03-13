# Description
It's a tree walk interpreter of [Lox language](https://craftinginterpreters.com/contents.html)
# Usage
## Start the REPL
`make run`
## Interpret a file
`make build && ./loxgo [file]`
## Produce Expression types
`make astgen && ./astgenerator pkg/parser/ast`