define i32 @f() {
0:
	%x = alloca i32
	%1 = load i32, i32* %x
	%2 = sub i32 0, %1
	ret i32 %2
}
define i32 @g() {
0:
	%y = alloca i32
	%1 = load i32, i32* %y
	%2 = icmp ne i32 %1, 0
	%3 = xor i1 %2, true
	%4 = zext i1 %3 to i32
	ret i32 %4
}
