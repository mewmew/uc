define i32 @f() {
0:
	%x = alloca i32
	%1 = load i32, i32* %x
	ret i32 %1
}
