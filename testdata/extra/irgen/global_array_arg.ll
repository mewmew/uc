@a = global [5 x i32] zeroinitializer

define void @g(i32* %a) {
0:
	%1 = alloca i32*
	store i32* %a, i32** %1
	ret void
}

define void @f() {
0:
	call void @g(i32* getelementptr ([5 x i32], [5 x i32]* @a, i64 0, i64 0))
	ret void
}
