gen:
	gocc -a uc.bnf
	find . -type f -name '*.go' | xargs goimports -w

debug:
	gocc -debug_lexer -v -a uc.bnf
	find . -type f -name '*.go' | xargs goimports -w

clean:
	rm -f errors/errors.go
	rm -f lexer/acttab.go
	rm -f lexer/lexer.go
	rm -f lexer/transitiontable.go
	rm -f parser/action.go
	rm -f parser/actiontable.go
	rm -f parser/gototable.go
	rm -f parser/parser.go
	rm -f parser/productionstable.go
	rm -f token/token.go
	rm -f util/litconv.go
	rm -f util/rune.go
	-rmdir errors
	-rmdir lexer
	-rmdir parser
	-rmdir token
	-rmdir util

.PHONY: gen debug clean
