#!/bin/bash
if [ $# -lt 1 ]; then
	echo "Usage: strip FILE.ll"
	exit 1
fi
f=$1

sar -i "target[^\n]+\n" "" "${f}"
sar -i "; <label>:([0-9]+)[^\n]+\n" "\$1:\n" "${f}"
sar -i ";[^\n]+\n" "" "${f}"
sar -i "attributes[^\n]+\n" "" "${f}"
sar -i " #[0-9]+" "" "${f}"
sar -i "[!][^\n]+\n" "" "${f}"
sar -i ", align [0-9]+\n" "\n" "${f}"
sar -i " nsw " " " "${f}"
sar -i " inbounds " " " "${f}"
sar -i "^\n" "" "${f}"
sar -i "\n\n" "\n" "${f}"
sar -i "\n\n" "\n" "${f}"
sar -i "  " "\t" "${f}"
sar -i "define[^\n]*" "\$0\n0:" "${f}"
sar -i " common " " " "${f}"
sar -i " signext " " " "${f}"
