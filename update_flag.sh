#!/bin/bash

clean=$(git status | grep "nothing to commit, working tree clean")
if [ -z "$clean" ]; then
    echo There are uncommitted changes.
    echo Updating will overwrite, commit or stash first.
    exit 1
else
    echo You might want to check if there are new files upstream:
    echo https://github.com/golang/go/tree/master/src/flag
    echo 
    curl -LO https://raw.githubusercontent.com/golang/go/master/src/flag/example_func_test.go
	curl -LO https://raw.githubusercontent.com/golang/go/master/src/flag/example_textvar_test.go
	curl -LO https://raw.githubusercontent.com/golang/go/master/src/flag/example_test.go
	curl -LO https://raw.githubusercontent.com/golang/go/master/src/flag/example_value_test.go
	curl -LO https://raw.githubusercontent.com/golang/go/master/src/flag/export_test.go
	curl -LO https://raw.githubusercontent.com/golang/go/master/src/flag/flag.go
	curl -LO https://raw.githubusercontent.com/golang/go/master/src/flag/flag_test.go

fi

