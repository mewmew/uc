define i32 @g(i8 %a, i8 %b) {
0:
	%1 = alloca i8
	%2 = alloca i8
	store i8 %a, i8* %1
	store i8 %b, i8* %2
	ret i32 42
}

define i32 @f() {
0:
	%x = alloca i32
	%y = alloca i32
	%1 = load i32, i32* %y
	%2 = trunc i32 %1 to i8
	%3 = load i32, i32* %x
	%4 = trunc i32 %3 to i8
	%5 = call i32 @g(i8 %2, i8 %4)
	ret i32 %5
}
