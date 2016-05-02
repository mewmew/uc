// Missing function return.
//
//    missing return at end of non-void function "f"
//
// Undefined behaviour
//    clang returns 13
//    gcc return 17
int f(int a) {
	;
}

int g(int a) {
	return a;
}

int main(void) {
	int a;
	a = 13;
	g(a);
	a = 17;
	return f(a);
}
