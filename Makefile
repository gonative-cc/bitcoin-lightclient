
setup:
	@cd .git/hooks; ln -s -f ../../scripts/git-hooks/* ./


.git/hooks/pre-commit: setup

build: .git/hooks/pre-commit
	go build .

start:
	@./btclightclient

clean:
	rm ./btclightclient

# used as pre-commit
lint-git:
	@files=$$(git diff --name-only --cached | grep  -E '\.go$$' | xargs -r gofmt -l); if [ -n "$$files" ]; then echo $$files;  exit 101; fi
	@git diff --name-only --cached | grep  -E '\.go$$' | xargs -r revive
	@git diff --name-only --cached | grep  -E '\.md$$' | xargs -r markdownlint-cli2

# lint changed files
lint:
	@files=$$(git diff --name-only | grep  -E '\.go$$' | xargs -r gofmt -l); if [ -n "$$files" ]; then echo $$files;  exit 101; fi
	@git diff --name-only | grep  -E '\.go$$' | xargs -r revive
	@git diff --name-only | grep  -E '\.md$$' | xargs -r markdownlint-cli2

lint-all: lint-fix-go-all
	@revive ./...

lint-fix-all: lint-fix-go-all

lint-fix-go-all:
	@gofmt -w -s -l .


.PHONY: build start clean setup
.PHONY: lint lint-all lint-fix-all lint-fix-go-all

###############################################################################
##                                   Tests                                   ##
###############################################################################

TEST_COVERAGE_PROFILE=coverage.txt
TEST_TARGETS := test-unit test-unit-cover test-race
test-unit: ARGS=-timeout=10m -tags='$(UNIT_TEST_TAGS)'
test-unit-cover: ARGS=-timeout=10m -tags='$(UNIT_TEST_TAGS)' -coverprofile=$(TEST_COVERAGE_PROFILE) -covermode=atomic
test-race: ARGS=-timeout=10m -race -tags='$(TEST_RACE_TAGS)'
$(TEST_TARGETS): run-tests

run-tests:
ifneq (,$(shell which tparse 2>/dev/null))
	@go test -mod=readonly -json $(ARGS) ./... | tparse
else
	@go test -mod=readonly $(ARGS) ./...
endif

cover-html: test-unit-cover
	@echo "--> Opening in the browser"
	@go tool cover -html=$(TEST_COVERAGE_PROFILE)
