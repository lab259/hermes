
GOPATH=$(CURDIR)/.gopath
GOPATHCMD=GOPATH=$(GOPATH)
PROJECT=github.com/jamillosantos/http
PROJECT_SRC=$(CURDIR)/src/$(PROJECT)
DEP=dep

COVERDIR=$(CURDIR)/.cover
COVERAGEFILE=$(COVERDIR)/cover.out

DEPS=$(call external_deps, '.')

.PHONY: get test test-watch coverage coverage-html

test:
	@${GOPATHCMD} ginkgo --failFast ./...

test-watch:
	@${GOPATHCMD} ginkgo watch -cover -r ./...

coverage:
	@mkdir -p $(COVERDIR)
	@${GOPATHCMD} ginkgo -r -covermode=count --cover --trace ./
	@echo "mode: count" > "${COVERAGEFILE}"
	@find . -type f -name *.coverprofile -exec grep -h -v "^mode:" {} >> "${COVERAGEFILE}" \; -exec rm -f {} \;

coverage-html:
	@$(GOPATHCMD) go tool cover -html="${COVERAGEFILE}" -o .cover/report.html

deps:
	@mkdir -p ${GOPATH}
	@go list -f '{{join .Deps "\n"}}' $(1) | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}' | GOPATH=${GOPATH} xargs go get
	@go list -f '{{join .TestImports "\n"}}' $(1) | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}' | GOPATH=${GOPATH} xargs go get

list-external-deps:
	$(call external_deps,'.')

define external_deps
	@echo '-- $(1)'; go list -f '{{join .Deps .TestImports " "}}' $(1) | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}'
endef
