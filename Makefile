all: pp-contact-tracer.pdf

%.pdf: %.md
	pandoc --standalone --table-of-contents --number-sections \
         --variable papersize=a4paper \
         --variable classoption=twocolumn \
         --variable links-as-notes \
         -s $< \
         -o $@

.PHONY: fmt clean
fmt:
	pandoc -o tmp.md -s pp-contact-tracer.md
	mv tmp.md pp-contact-tracer.md
	pandoc -o tmp.md -s README.md
	mv tmp.md README.md

clean:
	rm -f pp-contact-tracer.pdf
