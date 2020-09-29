package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/otoolep/json-vs-proto/chinook"
	"github.com/otoolep/json-vs-proto/command"
)

type POD struct {
	Transaction bool            `json:"tx,omitempty"`
	Timings     bool            `json:"timings,omitempty"`
	Query       []string        `json:"query,omitempty"`
	Value       [][]interface{} `json:"value,omitempty`
}

func main() {
	fmt.Println("Compare sizes....")

	//q1 := []string{"INSERT INTO foo(name, age, place, booked) VALUES(fiona, 20, GLENMORE, true)"}
	//q1 := []string{"SELECT * FROM foo", "INSERT INTO foo(name, age) VALUES(fiona, 20)", "INSERT INTO foo(name, age) VALUES(dana, 44)"}
	//q1 := []string{"INSERT INTO foo(name, age) VALUES(fiona, 20", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo"}
	q1 := []string{chinook.DB}

	p := &command.QueryCommand{
		Timings:     true,
		Transaction: false,
		Query:       q1,
		Value:       nil,
	}

	pb, err := proto.Marshal(p)
	if err != nil {
		log.Fatalf("failed to marshal protobuf: %s", err.Error())
	}
	pl := len(pb)
	fmt.Println("Proto:", pl)

	j := &POD{
		Timings:     true,
		Transaction: false,
		Query:       q1,
		Value:       nil,
	}

	jb, err := json.Marshal(j)
	if err != nil {
		log.Fatalf("failed to marshal JSON: %s", err.Error())
	}
	jl := len(jb)
	fmt.Println("JSON:", jl)

	fmt.Println(float64((jl-pl))/float64(jl)*100, "% reduction.")

	var buf bytes.Buffer
	zw, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		log.Fatal(err)
	}

	_, err = zw.Write(pb)
	if err != nil {
		log.Fatal(err)
	}
	zw.Close()
	fmt.Println("Buffer length:", buf.Len())
}
