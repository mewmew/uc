/* Test file for semantic errors. Contains exactly one error. */

int d(int a, int b);

int main (void){	
  d(1, 2, 3);	// Too many arguments to function 'd'
}
