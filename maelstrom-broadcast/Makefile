.gomod:
	go mod tidy

.gobuild: .gomod
	go install maelstrom-broadcast

.broadcast-single:
	${MAELSTROM_BIN} test -w broadcast --bin ${OUR_GO_PATH}/bin/maelstrom-broadcast --node-count 1 --time-limit 20 --rate 10

default: .gobuild .broadcast-single
