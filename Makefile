VERSION = $(shell cat VERSION)
VERSION_MAJOR = $(shell cat VERSION | sed 's/^v\([0-9]*\)\.\([0-9]*\)\.\([0-9]*\).*$$/\1/')
VERSION_MINOR = $(shell cat VERSION | sed 's/^v\([0-9]*\)\.\([0-9]*\)\.\([0-9]*\).*$$/\2/')
VERSION_BUILD = $(shell cat VERSION | sed 's/^v\([0-9]*\)\.\([0-9]*\)\.\([0-9]*\).*$$/\3/')
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
.PHONY: test test-ci


bump: VERSION
ifneq ($(BRANCH),master)
	echo "Bump only works on master, currently ${BRANCH}"
else
	$(eval VERSION_BUILD=$(shell echo "${VERSION_BUILD}+1"|bc))
	$(eval VERSION=v$(shell echo "${VERSION_MAJOR}.${VERSION_MINOR}.${VERSION_BUILD}"))
	echo "${VERSION}" > VERSION
	git add VERSION
	git commit -m "Version bump ${VERSION}"
	git tag -m "Tag ${VERSION}" ${VERSION}
	git push origin master
	git push --tags
endif
test:
	go test -v -cpu 1,4 $(shell glide novendor)

test-ci:
	go test -v -covermode=count  -bench . -cpu 1,4 $(shell glide novendor)
