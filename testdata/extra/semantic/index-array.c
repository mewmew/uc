// Invalid index expression
//
//    invalid array index; expected integer, got "int[20]"
void f(void) {
	int x[20];
	int y[20];
	x[y];
}
