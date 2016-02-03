/* Test file for syntactic errors. Contains exactly one error. */

int a;
int main(void) {
  while (a < 0) 	// Missing statement
}
