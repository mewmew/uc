define void @f(i32* %a) {
0:
	%1 = alloca i32*
	%b = alloca [5 x i32]
	store i32* %a, i32** %1
	%2 = getelementptr [5 x i32], [5 x i32]* %b, i32 0, i32 0
	store i32* %2, i32** %1
	ret void
}
