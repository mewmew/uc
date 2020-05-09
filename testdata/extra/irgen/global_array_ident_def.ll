@b = global [5 x i32] zeroinitializer

define void @f(i32* %a) {
0:
	%1 = alloca i32*
	store i32* %a, i32** %1
	store i32* getelementptr ([5 x i32], [5 x i32]* @b, i64 0, i64 0), i32** %1
	ret void
}
