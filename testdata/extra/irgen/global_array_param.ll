@b = global [5 x i32] zeroinitializer

define i32 @f(i32* %a) {
0:
	%1 = alloca i32*
	store i32* %a, i32** %1
	%2 = load i32*, i32** %1
	%3 = getelementptr i32, i32* %2, i64 6
	%4 = load i32, i32* %3
	ret i32 %4
}

define void @g() {
0:
	%1 = call i32 @f(i32* getelementptr ([5 x i32], [5 x i32]* @b, i64 0, i64 0))
	ret void
}
