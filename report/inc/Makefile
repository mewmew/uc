.SECONDARY:
TESTDATA=../../uc/testdata
TESTOUTDIR=testoutput
VPATH=$(TESTDATA):$(TESTOUTDIR)

all: clean 2_lexical_analysis

LEXSOURCEDIRS=$(wildcard $(TESTDATA)/*/lexer)
LEXSOURCE=$(wildcard $(LEXSOURCEDIRS:=/*.c))
LEXOBJECTDIRS=$(LEXSOURCEDIRS:$(TESTDATA)/%=$(TESTOUTDIR)/%)
LEXOBJECTS=$(LEXSOURCE:$(TESTDATA)/%.c=$(TESTOUTDIR)/%.txt)

UCFLAGS=
REDIRECTS=2>&1

$(LEXOBJECTDIRS):
	-mkdir -p $@

$(LEXOBJECTS): | $(LEXOBJECTDIRS)

2_lexical_analysis: $(LEXOBJECTDIRS:%=%/listing.tex)

%/listing.tex: $(LEXOBJECTS) #%/*.txt  #$(filter %,$(LEXOBJECTS))
	ls $(@D)/*.txt | awk '{printf "\\lstinputlisting[style=go]{inc/%s}\n\n", $$1}' >> $@


$(TESTOUTDIR)/incorrect/lexer/ugly.txt: UCFLAGS+=-n 10

$(TESTOUTDIR)/incorrect/lexer/%.txt: REDIRECTS+=1>/dev/null

$(TESTOUTDIR)/%.txt: $(TESTDATA)/%.c | $(LEXOBJECTDIRS)
	cd $(CURDIR)/$(TESTDATA); ulex $(UCFLAGS) $(<:$(TESTDATA)/%=%) > $(CURDIR:$(TESTDATA)/%=%)/$@ $(REDIRECTS)

clean:
	rm -rf $(TESTOUTDIR)

.PHONY: all 2_lexical_analysis clean $(LEXOBJECTDIRS)
