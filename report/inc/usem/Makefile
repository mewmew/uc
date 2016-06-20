TESTDATA=../../../testdata
UPROG=usem
UCFLAGS=-no-colors
REDIRECTS=2>&1

all: semantic_analysis

vpath %.c $(TESTDATA)
SOURCES=quiet/semantic/s01.c \
quiet/semantic/s02.c \
quiet/semantic/s03.c \
quiet/semantic/s04.c \
quiet/semantic/s05.c \
quiet/semantic/s06.c \
incorrect/semantic/se01.c \
incorrect/semantic/se02.c \
incorrect/semantic/se03.c \
incorrect/semantic/se04.c \
incorrect/semantic/se05.c \
incorrect/semantic/se06.c \
incorrect/semantic/se07.c \
incorrect/semantic/se08.c \
incorrect/semantic/se09.c \
incorrect/semantic/se10.c \
incorrect/semantic/se11.c \
incorrect/semantic/se12.c \
incorrect/semantic/se13.c \
incorrect/semantic/se14.c \
incorrect/semantic/se15.c \
incorrect/semantic/se16.c \
incorrect/semantic/se17.c \
incorrect/semantic/se18.c \
incorrect/semantic/se19.c \
incorrect/semantic/se20.c \
incorrect/semantic/se21.c \
incorrect/semantic/se22.c \
incorrect/semantic/se23.c \
incorrect/semantic/se24.c \
incorrect/semantic/se25.c \
incorrect/semantic/se26.c \
incorrect/semantic/se27.c \
incorrect/semantic/se28.c \
incorrect/semantic/se29.c \
incorrect/semantic/se30.c \
incorrect/semantic/se31.c \
incorrect/semantic/se32.c \
incorrect/semantic/se33.c \
incorrect/semantic/se34.c \
noisy/advanced/8queens.c \
noisy/advanced/bubble.c \
noisy/advanced/eval.c \
noisy/advanced/primes.c \
noisy/advanced/quick.c \
extra/semantic/extra-void-arg.c \
extra/semantic/incompatible-arg-type.c \
extra/semantic/index-array.c \
extra/semantic/local-var-redef.c \
extra/semantic/missing-return.c \
extra/semantic/missing-return-main.c \
extra/semantic/nested-function-def.c \
extra/semantic/param-redef.c \
extra/semantic/tentative-var-def.c \
extra/semantic/unnamed-arg.c \
extra/semantic/variable-sized-array.c \
extra/semantic/variable-sized-array-arg.c \
extra/semantic/void-array.c \
extra/semantic/void-array-arg.c \
extra/semantic/void-param.c \
extra/semantic/void-params.c \
extra/semantic/void-var.c

#$(patsubst $(TESTDATA)/%,%,$(wildcard $(TESTDATA)/noisy/advanced/*.c))

OBJECTS=$(SOURCES:%.c=%.txt)
OBJECTDIRS=$(sort $(dir $(OBJECTS)))

UPROGPATH=$(shell which $(UPROG))

semantic_analysis: $(OBJECTDIRS:=listing.tex)

$(OBJECTDIRS):
	mkdir -p "$@"

$(OBJECTS): %.txt: %.c $(UPROGPATH) | $(OBJECTDIRS)
	@echo making $@
	@cd $(TESTDATA) && $(UPROG) $(UCFLAGS) $${OLDPWD}/$< > $${OLDPWD}/$@ $(REDIRECTS)

# Does not work for adding of older files, only files newer than Makefile.depend
Makefile.depend: $(OBJECTS)
	@echo making $@
	@rm -f $@
	@(for a in $(OBJECTDIRS); do echo -n $$a"listing.tex: " >>$@; b=` for i in $(OBJECTS); do echo $$i; done | grep $$a `; echo $$b >>  $@; done;)

# Gets its prerequisits from Makefile.depend
# Also makes listref.tex
%listing.tex:
	@echo making $@ $(@D)/listref.tex
	@-rm -f $@ $(@D)/listref.tex
	@for txt in $^; do\
		c="$${txt%%.txt}.c";\
		printf "\\\\lstinputlisting[style=go,\n\
		caption=%s,\n\
		label=lst:$(UPROG)/%s\n\
		]{inc/$(UPROG)/%s}\n\n" "$$c" "$$c" "$$txt" >> $@;\
		printf "\n\nListing~\\\\ref{lst:$(UPROG)/%s}:~\\\\nameref{lst:$(UPROG)/%s}"\
		"$$c" "$$c" >> $(@D)/listref.tex; \
	done

clean:
	rm -rf $(OBJECTDIRS) Makefile.depend

.PHONY: all semantic_analysis clean

ifneq (clean,$(MAKECMDGOALS))
-include Makefile.depend
endif
