/* Test file for semantic errors. Contains exactly one error. */

char a[10];

int main (void) {
  a+1;	// Attempt to apply arithmetic to array reference
        // This is legal in C, but not in uC
}
