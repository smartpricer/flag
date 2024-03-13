include variables.mk
include go.mk

.PHONY: update

all: update

update:
	bash update_flag.sh

tag-and-push:
	git tag $(cat VERSION)
	git push origin $(cat VERSION)