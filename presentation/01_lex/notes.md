# Lexer

## Examples

### Correct

* testdata/foo/bar.c

### Incorrect

* testdata/foo/baz.c

## What is interesting with the design and implementation

* Hand-written
    - State function, Rob Pike.
    - Encoding

* Gocc generated
    - From BNF grammar

* The Lexer Hack!!
    - Type keywords are identifiers

## What was difficult?

* Line-comments ending with EOF
    - Hand-written lexer solved 
        + EOF case initially not handled, now using two cases, one for new line one for EOF.
    - Gocc generated lexer solution:
        + Insert new-line at end of input.
