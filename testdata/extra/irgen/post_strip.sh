#!/bin/bash
for f in "array_ident_use.ll"; do
	sar -i "getelementptr ([^,]+), ([^,]+), i32 0, i32 0" "getelementptr \$1, \$2, i64 0, i64 0" "${f}"
done
