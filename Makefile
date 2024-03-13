include variables.mk
include go.mk

.PHONY: update

all: update

update:
	bash update_flag.sh

