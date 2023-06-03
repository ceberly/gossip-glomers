export MAELSTROM_BIN=~/local/bin/maelstrom/maelstrom
export OUR_GO_PATH=`go env GOPATH`

echo:
	cd maelstrom-echo && make

unique-ids:
	cd maelstrom-unique-ids && make

broadcast:
	cd maelstrom-broadcast && make
