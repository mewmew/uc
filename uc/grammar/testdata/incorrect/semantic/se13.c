/* Test file for semantic errors. Contains exactly one error. */

void foo(int n){
  n;
}

int main(void) {
  1 + foo(0);	// 'foo' does not return a value
}
