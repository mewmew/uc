define void @f() {
0:
	%a = alloca [5 x i32]
	%1 = getelementptr [5 x i32], [5 x i32]* %a, i32 0, i32 0
	ret void
}
