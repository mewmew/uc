@a = global [5 x i32] zeroinitializer
define void @g(i32* %a) {
; <label>:0
	%1 = alloca i32*
	store i32* %a, i32** %1
	ret void
}
define void @f() {
; <label>:0
	%1 = getelementptr [5 x i32], [5 x i32]* @a, i64 0, i64 0
	call void @g(i32* %1)
	ret void
}
