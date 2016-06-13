# *.c -> *.ll
CFILES = $(wildcard *.c)
LLFILES = $(CFILES:.c=.ll)

all: $(LLFILES) post_strip

%.ll: %.c
	clang -S -emit-llvm -o $@ $<
	./strip.sh $@

post_strip:
	./post_strip.sh

.PHONY: post_strip
