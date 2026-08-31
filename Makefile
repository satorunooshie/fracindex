GO ?= go
FUZZTIME ?= 30s

.PHONY: test
test:
	$(GO) test ./...

.PHONY: vet
vet:
	$(GO) vet ./...

.PHONY: fuzz-smoke
fuzz-smoke:
	$(GO) test -run='^Fuzz' ./...

.PHONY: bench-smoke
bench-smoke:
	$(GO) test -run='^$$' -bench=. -benchtime=1x ./...

.PHONY: check
check: test vet fuzz-smoke bench-smoke

.PHONY: fuzz
fuzz:
	$(GO) test -fuzz=FuzzIndexerOperationSequence -fuzztime=$(FUZZTIME) ./...
	$(GO) test -fuzz=FuzzCustomAlphabetOperationSequence -fuzztime=$(FUZZTIME) ./...

.PHONY: bench
bench:
	$(GO) test -bench=. -benchmem ./...
