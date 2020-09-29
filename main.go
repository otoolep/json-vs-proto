package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	//"github.com/otoolep/json-vs-proto/chinook"
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

	//q1 := []string{"SELECT column1, column2 FROM table1 WHERE column3 IN ( SELECT TOP(1) column4 FROM table2 INNER JOIN table3 ON table2.column1 = table3.column1)"}
	//q1 := []string{"INSERT INTO foo(name, age, place, booked) VALUES(fiona, 20, GLENMORE, true)"}
	//q1 := []string{"INSERT INTO foo(name, age, place, booked) VALUES(fiona, 20, GLENMORE, true)", "SELECT column1, column2 FROM table1 WHERE column3 IN ( SELECT TOP(1) column4 FROM table2 INNER JOIN table3 ON table2.column1 = table3.column1)"}
	//q1 := []string{"SELECT * FROM foo", "INSERT INTO foo(name, age) VALUES(fiona, 20)", "INSERT INTO foo(name, age) VALUES(dana, 44)"}
	//q1 := []string{"INSERT INTO foo(name, age) VALUES(fiona, 20", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo"}
	q1 := []string{"SELECT * FROM foo", "SELECT column1, column2 FROM table1 WHERE column3 IN ( SELECT TOP(1) column4 FROM table2 INNER JOIN table3 ON table2.column1 = table3.column1)", "INSERT INTO foo(name, age, place, booked) VALUES(fiona, 20, GLENMORE, true)"}
	//q1 := []string{chinook.DB}

	// Proto only.
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

	// JSON only.
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

	// Compressed JSON
	var buf bytes.Buffer
	gz, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewEncoder(gz).Encode(j); err != nil {
		log.Fatalf("failed to encode and gzip JSON: %s", err.Error())
	}
	gz.Close()
	lbuf := buf.Len()
	fmt.Println("Compressed JSON:", lbuf)

	fmt.Print(int(float64((jl-pl))/float64(jl)*100), "% reduction moving from JSON to Proto.\n")
	fmt.Print(int(float64((jl-lbuf))/float64(jl)*100), "% reduction moving from JSON to compressed JSON.\n")

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
