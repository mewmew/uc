@b = global [5 x i32] zeroinitializer
define i32 @f(i32* %a) {
; <label>:0
	%1 = alloca i32*
	store i32* %a, i32** %1
	%2 = load i32*, i32** %1
	%3 = getelementptr i32, i32* %2, i64 6
	%4 = load i32, i32* %3
	ret i32 %4
}
define i32 @main() {
; <label>:0
	%1 = alloca i32
	store i32 0, i32* %1
	%2 = call i32 @f(i32* getelementptr ([5 x i32], [5 x i32]* @b, i32 0, i32 0))
	ret i32 %2
}
