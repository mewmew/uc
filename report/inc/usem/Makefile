.SECONDARY:
TESTDATA=../../../testdata
VPATH=$(CURDIR):$(TESTDATA)
UPROG=usem

all: clean semantic_analysis

SOURCEDIRS=$(wildcard $(TESTDATA)/*/semantic)
SOURCEDIRS+=$(wildcard $(TESTDATA)/noisy/*)
SOURCE=$(wildcard $(SOURCEDIRS:=/*.c))
OBJECTDIRS=$(SOURCEDIRS:$(TESTDATA)/%=%)
OBJECTS=$(SOURCE:$(TESTDATA)/%.c=%.txt)

UCFLAGS=
REDIRECTS=2>&1

semantic_analysis: $(OBJECTDIRS:=/listing.tex)

$(OBJECTDIRS):
	mkdir -p "$@"

$(OBJECTS): %.txt: $(TESTDATA)/%.c | $(OBJECTDIRS)
	cd $(CURDIR)/$(TESTDATA); $(UPROG) $(UCFLAGS) $(<:$(TESTDATA)/%=%) > $(CURDIR:$(TESTDATA)/%=%)/$@ $(REDIRECTS)

%/listing.tex: $(OBJECTS)
	ls $(@D)/*.txt | awk '{name=$$1;\
		sub(/.txt/,".c",name); printf "\\lstinputlisting[style=go,\
		caption=%s,label=lst:$(UPROG)/%s]{inc/$(UPROG)/%s}\n\n", name, name, $$1}' > $@
	ls $(@D)/*.txt | awk '{name=$$1;\
		sub(/.txt/,".c",name); printf "\nListing~\\ref{lst:$(UPROG)/%s}:~\\nameref{lst:$(UPROG)/%s}\n",\
		name, name}' > $(@D)/listref.tex

incorrect/er/ugly.txt: UCFLAGS+=-n 10

incorrect/er/%.txt: REDIRECTS+=1>/dev/null

clean:
	rm -rf $(OBJECTDIRS)

.PHONY: all semantic_analysis $(OBJECTDIRS) clean 
