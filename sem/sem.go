// Package sem implements a set of semantic analysis passes.
package sem

// TODO: Verify that all declarations occur at the beginning of the function
// body, and after the first non-declaration statement, no other declarations
// should be allowed to occur. Note, this pass should only be enabled for older
// versions of C, as newer ones allow declarations to occur throughout the
// function (albeit with other restrictions, e.g. goto may not jump over
// declarations).
