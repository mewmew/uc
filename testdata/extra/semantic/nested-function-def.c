// Nested functions valid in gcc and our implementation, but invalid in uC
//
// Activate check for nested functions with the -no-nested-functions flag
int main(void) {
	void f(void){
		;
	}
}
