#!/bin/bash
if [ $# -lt 1 ]; then
	echo "Usage: compile FILE.c"
	exit 1
fi
f=$1
uclang -o a.ll "${f}"
if [ $? -ne 0 ]; then
	echo "FAILURE: ${f}"
	exit 1
fi
llvm-link -S -o out.ll a.ll testdata/uc.ll
echo "SUCCESS: ${f}"
