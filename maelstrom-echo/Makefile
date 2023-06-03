.gomod:
	go mod tidy

.gobuild: .gomod
	go install maelstrom-echo

default: .gobuild
	${MAELSTROM_BIN} test -w echo --bin ${OUR_GO_PATH}/bin/maelstrom-echo --node-count 1 --time-limit 10
