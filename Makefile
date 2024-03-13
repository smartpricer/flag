include variables.mk
include go.mk

.PHONY: update

all: update

update:
	./update_flag.sh

