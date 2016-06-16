@a = global i32 0
define i32 @f() {
0:
	%1 = load i32, i32* @a
	ret i32 %1
}
