PROC ?= $(shell sh scripts/print_proc.sh)
BIN ?= serv
HASH := $(shell sh -c 'git rev-parse --verify HEAD | cut -c 1-8')

gen:
	bash ./scripts/mock_gen.sh
	