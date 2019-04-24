GOPATH=$(CURDIR)/../../../../
GOPATHCMD=GOPATH=$(GOPATH)

COVERDIR=$(CURDIR)/.cover
COVERAGEFILE=$(COVERDIR)/cover.out

.PHONY: dep-ensure dep-update vet test test-watch coverage coverage-ci coverage-html

dep-ensure:
	@$(GOPATHCMD) dep ensure -v

dep-update:
	@$(GOPATHCMD) dep update -v $(PACKAGE)

vet:
	@$(GOPATHCMD) go vet ./...

fmt:
	@$(GOPATHCMD) go fmt ./...

test:
	@${GOPATHCMD} ginkgo --failFast ./...

test-watch:
	@${GOPATHCMD} ginkgo watch -cover -r ./...

coverage: coverage-ci
	@sed -i -e "s|_$(CURDIR)/|./|g" "${COVERAGEFILE}"
	@cp "${COVERAGEFILE}" coverage.txt

coverage-ci:
	@mkdir -p $(COVERDIR)
	@${GOPATHCMD} ginkgo -r -covermode=count --cover --trace ./
	@echo "mode: count" > "${COVERAGEFILE}"
	@find . -type f -name *.coverprofile -exec grep -h -v "^mode:" {} >> "${COVERAGEFILE}" \; -exec rm -f {} \;

coverage-html:
	@$(GOPATHCMD) go tool cover -html="${COVERAGEFILE}" -o .cover/report.html

deps:
	@mkdir -p ${GOPATH}
	@go list -f '{{join .Deps "\n"}}' . | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}' | GOPATH=${GOPATH} xargs go get
	@go list -f '{{join .TestImports "\n"}}' . | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}' | GOPATH=${GOPATH} xargs go get

deps-ci:
	-go get -v -t ./...

list-external-deps:
	$(call external_deps,'.')

define external_deps
	@echo '-- $(1)'; go list -f '{{join .Deps .TestImports " "}}' $(1) | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}'
endef
