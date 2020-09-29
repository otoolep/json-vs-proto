# json-vs-proto

- Set GOBIN and add to PATH:
go get -u github.com/golang/protobuf/protoc-gen-go
go install github.com/golang/protobuf/protoc-gen-go

protoc -I=/home/philip/repos/json-vs-proto/src/github.com/otoolep/json-vs-proto/proto/ --go_out=/home/philip/repos/json-vs-proto/src/github.com/otoolep/json-vs-proto/proto/ /home/philip/repos/json-vs-proto/src/github.com/otoolep/json-vs-proto/proto/command.proto 
