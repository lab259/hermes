GOPATH=$(CURDIR)/../../../../
GOPATHCMD=GOPATH=$(GOPATH)

COVERDIR=$(CURDIR)/.cover
COVERAGEFILE=$(COVERDIR)/cover.out

EXAMPLES=$(shell ls ./examples/)

.PHONY: run dep-ensure dep-update vet test test-watch coverage coverage-ci coverage-html

run:
	@$(GOPATHCMD) go run examples/$(EXAMPLE)/main.go

build:
	@test -d ./examples && $(foreach example,$(EXAMPLES),$(GOPATHCMD) go build "-ldflags=$(LDFLAGS)" -o ./bin/$(example) -v ./examples/$(example) &&) :

dep-ensure:
	@$(GOPATHCMD) dep ensure -v

dep-update:
	@$(GOPATHCMD) dep update -v $(PACKAGE)

vet:
	@$(GOPATHCMD) go vet ./...

fmt:
	@$(GOPATHCMD) gofmt -e -s -d *.go

test:
	@${GOPATHCMD} ginkgo --failFast ./...

test-watch:
	@${GOPATHCMD} ginkgo watch -cover -r ./...

bench:
	@mkdir -p ./bench-results
	@${GOPATHCMD} go test -benchmem -run=github.com/lab259/http -bench=$(TARGET)$$ -test.parallel=1 -cpuprofile ./bench-results/profile.out
	@${GOPATHCMD} go tool pprof -svg profile.out > ./bench-results/profile.svg

coverage-ci:
	@mkdir -p $(COVERDIR)
	@${GOPATHCMD} ginkgo -r -covermode=count --cover --trace ./
	@echo "mode: count" > "${COVERAGEFILE}"
	@find . -type f -name *.coverprofile -exec grep -h -v "^mode:" {} >> "${COVERAGEFILE}" \; -exec rm -f {} \;

coverage: coverage-ci
	@sed -i -e "s|_$(CURDIR)/|./|g" "${COVERAGEFILE}"
	@cp "${COVERAGEFILE}" coverage.txt

coverage-html: coverage
	@$(GOPATHCMD) go tool cover -html="${COVERAGEFILE}" -o .cover/report.html
