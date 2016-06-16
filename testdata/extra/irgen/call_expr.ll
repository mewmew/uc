define i32 @g(i32 %a) {
0:
	%1 = alloca i32
	store i32 %a, i32* %1
	ret i32 42
}
define i32 @f() {
0:
	%1 = call i32 @g(i32 1)
	ret i32 %1
}
