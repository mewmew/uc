
all: ulex uparse usem

ulex uparse usem:
	make -C $@/

clean:
	make -C uparse clean
	make -C ulex clean
	make -C usem clean

.PHONY: all ulex uparse usem clean
