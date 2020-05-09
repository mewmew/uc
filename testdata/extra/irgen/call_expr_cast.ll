define i32 @g(i8 %a) {
0:
	%1 = alloca i8
	store i8 %a, i8* %1
	ret i32 42
}

define i32 @f() {
0:
	%x = alloca i32
	%1 = load i32, i32* %x
	%2 = trunc i32 %1 to i8
	%3 = call i32 @g(i8 %2)
	ret i32 %3
}
