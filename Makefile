VERSION = $(shell cat VERSION)
VERSION_MAJOR = $(shell cat VERSION | sed 's/^v\([0-9]*\)\.\([0-9]*\)\.\([0-9]*\).*$$/\1/')
VERSION_MINOR = $(shell cat VERSION | sed 's/^v\([0-9]*\)\.\([0-9]*\)\.\([0-9]*\).*$$/\2/')
VERSION_BUILD = $(shell cat VERSION | sed 's/^v\([0-9]*\)\.\([0-9]*\)\.\([0-9]*\).*$$/\3/')

.PHONY: test test-ci


bump: VERSION
	$(eval VERSION_BUILD=$(shell echo "${VERSION_BUILD}+1"|bc))
	$(eval VERSION=v$(shell echo "${VERSION_MAJOR}.${VERSION_MINOR}.${VERSION_BUILD}"))
	echo "${VERSION}" > VERSION
	git add VERSION
	git commit -m "Version bump ${VERSION}"
	git tag -m "Tag ${VERSION}" ${VERSION}
	git push origin master
	git push --tags

test:
	go test -v -cpu 1,4 $(shell glide novendor)

test-ci:
	go test -v -covermode=count  -bench . -cpu 1,4 $(shell glide novendor)
