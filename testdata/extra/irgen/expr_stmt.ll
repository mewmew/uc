define void @g() {
0:
	ret void
}

define void @f() {
0:
	call void @g()
	ret void
}
