define void @f(i32* %x) {
0:
	%1 = alloca i32*
	store i32* %x, i32** %1
	%2 = load i32*, i32** %1
	ret void
}
