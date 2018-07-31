
GOPATH=$(CURDIR)/.gopath
GOPATHCMD=GOPATH=$(GOPATH)
PROJECT=github.com/jamillosantos/http
PROJECT_SRC=$(CURDIR)/src/$(PROJECT)
DEP=dep

COVERDIR=$(CURDIR)/.cover
COVERAGEFILE=$(COVERDIR)/cover.out

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
	@GOPATH=${GOPATH} go get -v -t ./...
	@GOPATH=${GOPATH} go test -i ./...

list-external-deps:
	$(call external_deps,'.')

restore-import-paths:
	find . -name '*.go' -type f -execdir sed -i '' s%\"github.com/$(REPO_OWNER)/migration%\"github.com/mattes/migrate%g '{}' \;

rewrite-import-paths:
	find . -name '*.go' -type f -execdir sed -i '' s%\"github.com/jamillosantos/migration%\"github.com/$(REPO_OWNER)/migrate%g '{}' \;
