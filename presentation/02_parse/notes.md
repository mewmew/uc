# Parser

## Examples

### Correct

* testdata/foo/bar.c

### Incorrect

* testdata/foo/baz.c

## What is interesting with the design and implementation

* Gocc generated
    - From BNF grammar

* Wrapper around the hand-written lexer to unify the API of the two lexer implementations.
    + Since the Gocc generated parser is already expecting the Gocc lexer API.

* For the cannonical LR-1 shift-reduce problem (dangling else)
    - Use the -a flag of Gocc so automatically resolve the conflict through maximal-munch.

* Simplify and unify the grammar
    - Get rid of the Locals production rule, allow local variable definitions to occur anywhere in a function.

* Construct the AST using production actions. << astx.NewFuncDecl(foo, bar) >>
    + Show BNF grammar, and demonstrate how gocc is able to generate a lexer and parser from it.

* Associativity and precedence of binary operations
    - Directly from the Dragon book and the C spec.
    - Tree structure to define precedence.
        + Show example
    - Left or right side of production rule to define associativity.
        + Show example

* User-friendly error messages
    - PR to Gocc. Basically print the next token and the expected followset when unable to shift or reduce.

## What was difficult?

* Reduce-reduce conflicts introduced by The Lexer Hack.
    - Use a single production rule for type names and postpone type analysis to the semantic analysis stage.

* Evaluate merging TopLevelDecl with Decl to simplify grammar
    - Turns out, this was probably not the right design choice, as it made the implementation of the irgen more difficult. We were no longer able to differentiate between local and global variables with ease.
