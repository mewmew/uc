/* Test file for semantic errors. Contains exactly one error. */

int a;
int main(void) {
  a = foo(a);	// Function 'foo' not defined
}
