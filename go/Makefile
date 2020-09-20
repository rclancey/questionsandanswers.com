GOPATH = $(CURDIR)

SRC = $(shell find src -type f | grep '\.go')

DEPS = github.com/satori/uuid \
	github.com/mattn/go-sqlite3 \
	github.com/jmoiron/sqlx

DEPS_SRC = $(foreach dep, $(DEPS), src/$(dep))

export GOPATH

all: $(PRODUCT)

deps: $(DEPS_SRC)

.PHONY: deps

src/github.com/%:
	go get github.com/$*

$(PRODUCT): deps $(SRC)
	go build -o $(PRODUCT) $(PRODUCT)

clean:
	rm -f $(PRODUCT)

.PHONY: clean

distclean: clean
	rm -rf $(DEPS_SRC)
	rm -rf pkg

.PHONY: distclean
