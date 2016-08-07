@a = global [5 x i32] zeroinitializer
define void @f() {
; <label>:0
	%1 = load i32, i32* getelementptr ([5 x i32], [5 x i32]* @a, i64 0, i64 2)
	ret void
}
