include variables.mk
include go.mk

.PHONY: update

all: update

update:
	bash update_flag.sh

tag-and-push:
	bash git_tag_push.sh