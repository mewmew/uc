define void @f(i8 %a) {
0:
	%1 = alloca i8
	store i8 %a, i8* %1
	ret void
}

define void @g() {
0:
	%a = alloca i32
	%1 = load i32, i32* %a
	%2 = trunc i32 %1 to i8
	call void @f(i8 %2)
	ret void
}
