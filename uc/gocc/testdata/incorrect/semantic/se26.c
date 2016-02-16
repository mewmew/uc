/* Test file for semantic errors. Contains exactly one error. */

void f(int b[]) {
  b[0]=0;
}

int main(void) {
  char a[10];
  f(a);
}

