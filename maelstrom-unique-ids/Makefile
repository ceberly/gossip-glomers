.gomod:
	go mod tidy

.gobuild: .gomod
	go install maelstrom-unique-ids

default: .gobuild
	${MAELSTROM_BIN} test -w unique-ids --bin ${OUR_GO_PATH}/bin/maelstrom-unique-ids --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition
