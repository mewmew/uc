# Semantic Analysis

## Examples

### Correct

* testdata/foo/bar.c

### Incorrect

* testdata/foo/baz.c

## What is interesting with the design and implementation

* Design choice, forward declarations. Thus resolve identifiers in two passes. First file-scope declarations, then block scope-declarations.
    - Sem check: Step 1. Identifier resolution. Step 2: Type-checking. Step 3: Semantic analysis (always disallow nested functions, for now).

* Tree walking

## What was difficult?

* Tentative definitions.
    - Distinguish between declaration and definition.
        + Solution: Scope with `IsDef` method, which is different at file-scope and block-scope.
