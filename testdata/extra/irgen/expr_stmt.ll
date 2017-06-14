define void @g() {
; <label>:0
	ret void
}

define void @f() {
; <label>:0
	call void @g()
	ret void
}
