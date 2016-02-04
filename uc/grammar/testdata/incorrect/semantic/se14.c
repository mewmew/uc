/* Test file for semantic errors. Contains exactly one error. */

int f(int n) {
  return n / 2;
}


void p(int n) {
  int f;

  f = n * 2;
  f(n);		// 'f' refers only to the local variable
}


int n;

int main(void) {
  p(n);
}

