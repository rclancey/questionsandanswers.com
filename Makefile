PRODUCT = questionsandanswers
VERSION = 0.0.1
RELEASE = $(PRODUCT)-$(VERSION)
BUILDDIR = build/$(RELEASE)
BUILD = $(RELEASE).tar.gz

export PRODUCT

all: build/$(BUILD)

go:
	env GOPATH=$(CURDIR)/go $(MAKE) -C go

.PHONY: go

js:
	cd js && yarn install && yarn build

.PHONY: js

build: go js
	mkdir -p $(BUILDDIR)/bin $(BUILDDIR)/htdocs
	cp go/$(PRODUCT) $(BUILDDIR)/bin/
	cp bin/*.sh $(BUILDDIR)/bin/
	rsync -a js/build/ $(BUILDDIR)/htdocs/

.PHONY: build

build/$(BUILD): build
	cd build && tar -zcf $(BUILD) $(RELEASE)

clean:
	rm -rf js/build/
	env GOPATH=$(CURDIR)/go $(MAKE) -C go clean
	rm -rf build

.PHONY: clean

distclean: clean
	rm -rf js/node_modules
	env GOPATH=$(CURDIR)/go $(MAKE) -C go distclean

.PHONY: distclean
