// Missing return is valid for main.
//
// "reaching the `}` that terminates the main function returns a value of 0"
// (§5.1.2.2.3 in the C11 spec)
int main(void) {
}
