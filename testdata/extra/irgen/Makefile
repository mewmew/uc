# *.c -> *.ll
CFILES = $(wildcard *.c)
LLFILES = $(CFILES:.c=.ll)

all: sar $(LLFILES) post_strip

%.ll: %.c
	clang -S -emit-llvm -o $@ $<
	./strip.sh $@

post_strip:
	./post_strip.sh

sar:
	@if ! which $@ &> /dev/null ; then \
		echo "Please install the \"sar\" tool"; \
		echo; \
		echo "   go get github.com/mewkiz/cmd/sar"; \
		echo; \
		exit 1; \
	fi

.PHONY: sar post_strip
