.gomod:
	go mod tidy

.gobuild: .gomod
	go install maelstrom-broadcast

.broadcast-a:
	${MAELSTROM_BIN} test -w broadcast --bin ${OUR_GO_PATH}/bin/maelstrom-broadcast --node-count 1 --time-limit 20 --rate 10

.broadcast-b:
	${MAELSTROM_BIN} test -w broadcast --bin ${OUR_GO_PATH}/bin/maelstrom-broadcast --node-count 5 --time-limit 20 --rate 10

.broadcast-c:
	${MAELSTROM_BIN} test -w broadcast --bin ${OUR_GO_PATH}/bin/maelstrom-broadcast --node-count 5 --time-limit 20 --rate 10 --nemesis partition

#default: .gobuild .broadcast-a .broadcast-b .broadcast-c
default: .gobuild .broadcast-a
