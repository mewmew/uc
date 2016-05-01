// Invalid function call
//
//    calling "a" with incompatible argument type "int" to parameter of type "int[]".
int a(int b[]){
	return b[0];
}

int main(void) {
	int b;
	return a(b);
}
