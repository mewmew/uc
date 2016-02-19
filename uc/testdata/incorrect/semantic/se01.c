/* Test file for semantic errors. Contains exactly one error. */

int a;
int main(void) {	
	a = a + b;	// Variable 'b' not defined
}
