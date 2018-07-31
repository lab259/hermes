
GOPATH=$(CURDIR)
GOPATHCMD=GOPATH=$(CURDIR)
PROJECT=github.com/jamillosantos/http
PROJECT_SRC=$(CURDIR)/src/$(PROJECT)
DEP=cd $(PROJECT_SRC) && GOPATH=${GOPATH} dep

COVERDIR=$(CURDIR)/.cover
COVERAGEFILE=$(COVERDIR)/cover.out

.PHONY: dep-ensure dep-add dep-status test test-watch coverage coverage-html

dep-ensure:
	@$(DEP) ensure -v

dep-add:
ifdef PACKAGE
	@$(DEP) ensure -v -add $(PACKAGE)
else
	@echo "Usage: PACKAGE=<package url> make dep-add"
	@echo "The environment variable \`PACKAGE\` is not defined."
endif

dep-status:
	@$(DEP) status

test:
ifdef TARGET
	@echo Running tests from src/${PROJECT}/${TARGET} tests
	@${GOPATHCMD} ginkgo --failFast ./src/${PROJECT}/${TARGET}/...
else
	@${GOPATHCMD} ginkgo --failFast ./src/${PROJECT}/...
endif

race:
ifdef TARGET
	@echo Running tests from src/${PROJECT}/${TARGET} tests
	@${GOPATHCMD} ginkgo --failFast ./src/${PROJECT}/${TARGET}/...
else
	@${GOPATHCMD} ginkgo --failFast ./src/${PROJECT}/...
endif

test-watch:
ifdef TARGET
	@echo Watching src/${PROJECT}/${TARGET} tests
	@AURE_DIR=$(CURDIR) ${GOPATHCMD} ginkgo watch -cover -r ./src/${PROJECT}/${TARGET}
else
	@AURE_DIR=$(CURDIR) ${GOPATHCMD} ginkgo watch -cover -r ./src/${PROJECT}
endif

coverage:
	@mkdir -p $(COVERDIR)
	@@AURE_DIR=$(CURDIR) ${GOPATHCMD} ginkgo -r -covermode=count --cover --trace ./src/${PROJECT}
	@echo "mode: count" > "${COVERAGEFILE}"
	@find . -type f -name *.coverprofile -exec grep -h -v "^mode:" {} >> "${COVERAGEFILE}" \; -exec rm -f {} \;

coverage-html: coverage
	@$(GOPATHCMD) go tool cover -html="${COVERAGEFILE}" -o .cover/report.html

vet:
	@$(GOPATHCMD) go vet ./src/${PROJECT}