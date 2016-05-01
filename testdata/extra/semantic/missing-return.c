// Missing return; valid for main.
//
//    warning: control reaches end of non-void function
//
// Undefined behaviour
//    clang returns 106
//    gcc return 45
int f(int a) {
	int d =111;
	a;
	//d;
}

int d(int a) {
	return a;
}

int main(void){
	int a =3434;
	int b =45;
	int c =24;
	a=d(a);
	return f(b);
}
