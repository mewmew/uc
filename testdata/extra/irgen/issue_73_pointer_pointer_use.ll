define i32 @f(i32* %x) {
0:
	%1 = alloca i32*
	store i32* %x, i32** %1
	%2 = load i32*, i32** %1
	%3 = getelementptr i32, i32* %2, i64 0
	%4 = load i32, i32* %3
	ret i32 %4
}
