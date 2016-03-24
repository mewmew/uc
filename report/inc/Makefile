.SECONDARY:
TESTDATA=../../uc/testdata
TESTOUTDIR=testoutput
VPATH=$(TESTDATA):$(TESTOUTDIR)

all: 2_lexical_analysis

LEXSOURCEDIRS=$(wildcard $(TESTDATA)/*/lexer)
LEXSOURCE=$(wildcard $(LEXSOURCEDIRS:=/*.c))
LEXOBJECTDIRS=$(LEXSOURCEDIRS:$(TESTDATA)/%=$(TESTOUTDIR)/%)
LEXOBJECTS=$(LEXSOURCE:$(TESTDATA)/%.c=$(TESTOUTDIR)/%.txt)

$(LEXOBJECTDIRS):
	-mkdir -p $@

$(LEXOBJECTS): | $(LEXOBJECTDIRS)

2_lexical_analysis: $(LEXOBJECTDIRS:%=%/listing.tex)

%/listing.tex: $(LEXOBJECTS) #%/*.txt  #$(filter %,$(LEXOBJECTS))
	ls $(@D)/*.txt | awk '{printf "\\lstinputlisting[style=go]{inc/%s}\n\n", $$1}' >> $@

$(TESTOUTDIR)/%.txt: $(TESTDATA)/%.c | $(LEXOBJECTDIRS)
	cd $(CURDIR)/$(TESTDATA); ulex -n 10 $(<:$(TESTDATA)/%=%) > $(CURDIR:$(TESTDATA)/%=%)/$@ 2>&1

clean:
	rm -rf $(TESTOUTDIR)

.PHONY: all 2_lexical_analysis clean $(LEXOBJECTDIRS)
