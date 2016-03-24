.SECONDARY:
TESTDATA=../../uc/testdata
TESTOUTDIR=testoutput
VPATH=$(TESTDATA):$(TESTOUTDIR)

all: 2_lexical_analysis

LEXSOURCEDIRS=$(wildcard $(TESTDATA)/*/lexer)
LEXSOURCE=$(wildcard $(LEXSOURCEDIRS:=/*.c))
LEXOBJECTDIRS=$(LEXSOURCEDIRS:$(TESTDATA)/%=$(TESTOUTDIR)/%)
LEXOBJECTS=$(LEXSOURCE:$(TESTDATA)/%.c=$(TESTOUTDIR)/%.txt)

#%.tex: $(LEXOBJECTS)
#	ls $(LEXOBJECTS) | awk '{printf "\\lstinputlisting[style=go]{inc/%s}\n", $$1}' >> $@

$(LEXOBJECTDIRS):
	-mkdir -p $@

2_lexical_analysis: $(LEXSOURCEDIRS:$(TESTDATA)/%=$(TESTOUTDIR)/%/listing.tex)

%/listing.tex: | $(LEXOBJECTS) #%/*.txt  #$(filter %,$(LEXOBJECTS))
	ls $(@D)/*.txt | awk '{printf "\\lstinputlisting[style=go]{inc/%s}\n", $$1}' >> $@

$(TESTOUTDIR)/%.txt: $(TESTDATA)/%.c
	-mkdir -p $(@D)
	ulex -n 10 $< > $@ 2>&1
	#ls "$@" | awk '{printf "\\lstinputlisting[style=go]{inc/%s}\n", $$1}' >> $(TESTOUTDIR)/list.tex

clean:
	rm -rf $(TESTOUTDIR)

.PHONY: all 2_lexical_analysis clean $(LEXOBJECTDIRS)
