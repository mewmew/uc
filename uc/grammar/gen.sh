gocc -v -a grammar.bnf
find . -type f -name '*.go' | xargs goimports -w
