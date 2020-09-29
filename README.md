# json-vs-proto

- Set GOBIN and add to PATH: export PATH=$PATH:$GOBIN

go get -u github.com/golang/protobuf/protoc-gen-go

go install github.com/golang/protobuf/protoc-gen-go

export SRC_DIR=/home/philip/repos/rqlite/src/github.com/rqlite/rqlite/store/proto
export DEST_DIR=/home/philip/repos/rqlite/src/github.com/rqlite/rqlite/store/proto

protoc -I=$SRC_DIR --go_out=$DEST_DIR $SRC_DIR/command.proto

Findings:
- Proto is always a win, the smaller the query, the bigger the win (JSON is a bigger overhead in that case)
- Compression more likely to be a win the more statements (>5?) there are OR if one of the queries is very long (100 characters?).