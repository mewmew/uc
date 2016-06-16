define void @g(i32* %a) {
0:
	%1 = alloca i32*
	store i32* %a, i32** %1
	ret void
}
define void @f() {
0:
	%a = alloca [5 x i32]
	%1 = getelementptr [5 x i32], [5 x i32]* %a, i32 0, i32 0
	call void @g(i32* %1)
	ret void
}
