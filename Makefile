.PHONY: install

default: totalspaces

totalspaces: $(shell find . -name "*.go")
	go build

install: /usr/local/bin/ts

/usr/local/bin/ts: totalspaces
	cp $< $@
