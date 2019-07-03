COVERDIR=$(CURDIR)/.cover
COVERAGEFILE=$(COVERDIR)/cover.out

EXAMPLES=$(shell ls ./examples/)

$(EXAMPLES): %:
	$(eval EXAMPLE=$*)
	@:

run:
	@if [ ! -z "$(EXAMPLE)" ]; then \
		go run ./examples/$(EXAMPLE); \
	else \
		echo "Usage: make [$(EXAMPLES)] run"; \
		echo "The environment variable \`EXAMPLE\` is not defined."; \
	fi

build:
	@test -d ./examples && $(foreach example,$(EXAMPLES),go build "-ldflags=$(LDFLAGS)" -o ./bin/$(example) -v ./examples/$(example) &&) :

vet:
	@go vet ./...

fmt:
	@go fmt ./...

test:
	@ginkgo --failFast ./...

test-watch:
	@ginkgo watch -cover -r ./...

bench:
	@mkdir -p ./bench-results
	@go test -benchmem -run=github.com/lab259/hermes -bench=$(TARGET)$$ -test.parallel=1 -cpuprofile ./bench-results/cpu.prof -memprofile ./bench-results/mem.prof
	
plot-cpu:
	@go tool pprof -http :8080 ./bench-results/cpu.prof

plot-mem:
	@go tool pprof -alloc_space -http :8080 ./bench-results/mem.prof

coverage-ci:
	@mkdir -p $(COVERDIR)
	@ginkgo -r -covermode=count --cover --trace ./
	@echo "mode: count" > "${COVERAGEFILE}"
	@find . -type f -name *.coverprofile -exec grep -h -v "^mode:" {} >> "${COVERAGEFILE}" \; -exec rm -f {} \;

coverage: coverage-ci
	@sed -i -e "s|_$(CURDIR)/|./|g" "${COVERAGEFILE}"
	@cp "${COVERAGEFILE}" coverage.txt

coverage-html: coverage
	@go tool cover -html="${COVERAGEFILE}" -o .cover/report.html
	@xdg-open .cover/report.html 2> /dev/null > /dev/null

.PHONY: run build fmt vet test test-watch coverage coverage-ci coverage-html bench plot-cpu plot-mem