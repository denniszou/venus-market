git=$(subst -,.,$(shell git describe --always --match=NeVeRmAtCh --dirty 2>/dev/null || git rev-parse --short HEAD 2>/dev/null))

ldflags=-X=github.com/filecoin-project/venus-market.GitCommit=${git}
ifneq ($(strip $(LDFLAGS)),)
	ldflags+=-extldflags=$(LDFLAGS)
endif

GOFLAGS+=-ldflags="$(ldflags)"

build:
	rm -rf venus-market
	go build $(GOFLAGS) -o venus-market .
	./venus-market --version

deps:
	git submodule update --init
	./extern/filecoin-ffi/install-filcrypto

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run

test:
	rm -rf models/test_sqlite_db*
	go test -race ./...

