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
	//q1 := []string{`SELECT * FROM "foo"`, "SELECT * FROM foo", "SELECT * FROM bar", "SELECT * FROM qux"}
	q1 := []string{"INSERT INTO foo(name, ag) VALUES(fiona, 20", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo", "SELECT * FROM foo"}
	//q1 := []string{"SELECT * FROM foo", "SELECT column1, column2 FROM table1 WHERE column3 IN ( SELECT TOP(1) column4 FROM table2 INNER JOIN table3 ON table2.column1 = table3.column1)", "INSERT INTO foo(name, age, place, booked) VALUES(fiona, 20, GLENMORE, true)"}
	//q1 := []string{chinook.DB}

	v0 := []interface{}{1, "qux"}
	v1 := [][]interface{}{v0}

	// Proto only.
	pv1 := &command.Parameter{
		Value: &command.Parameter_I{
			I: 1,
		},
	}
	p := &command.QueryCommand{
		Timings:     true,
		Transaction: false,
		Query:       q1,
		Value:       []*command.Parameter{pv1},
	}
	pb, err := proto.Marshal(p)
	if err != nil {
		log.Fatalf("failed to marshal protobuf: %s", err.Error())
	}
	pl := len(pb)
	fmt.Println("Proto:", pl)

	// Compressed proto
	var cp bytes.Buffer
	gzp, err := gzip.NewWriterLevel(&cp, gzip.BestCompression)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewEncoder(gzp).Encode(pb); err != nil {
		log.Fatalf("failed to compress marshalled proto: %s", err.Error())
	}
	gzp.Close()
	lcp := cp.Len()
	fmt.Println("Compressed marshalled proto", lcp)

	// JSON only.
	j := &POD{
		Timings:     true,
		Transaction: false,
		Query:       q1,
		Value:       v1,
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

	// Encode SQL as JSON, compress, add to proto.
	var bufq bytes.Buffer
	gzq, err := gzip.NewWriterLevel(&bufq, gzip.BestCompression)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewEncoder(gzq).Encode(q1); err != nil {
		log.Fatalf("failed to encode as JSON and compress queries: %s", err.Error())
	}
	gzq.Close()
	pq := &command.QueryCommand{
		Timings:         true,
		Transaction:     false,
		CompressedQuery: bufq.Bytes(),
	}
	pbq, err := proto.Marshal(pq)
	if err != nil {
		log.Fatalf("failed to marshal protobuf: %s", err.Error())
	}
	jql := len(pbq)
	fmt.Println("Compressed JSON in proto:", jql)

	fmt.Print(int(float64((jl-pl))/float64(jl)*100), "% reduction moving from JSON to Proto.\n")
	fmt.Print(int(float64((jl-lcp))/float64(jl)*100), "% reduction moving from JSON to compressed marshalled Proto.\n")
	fmt.Print(int(float64((jl-lbuf))/float64(jl)*100), "% reduction moving from JSON to compressed JSON.\n")
	fmt.Print(int(float64((jl-jql))/float64(jl)*100), "% reduction moving from JSON to compressed JSON queries in Proto.\n")

	fmt.Println("===========================")
	fmt.Println("New model")
	fmt.Println("===========================")

        nmv := &command.Parameter{
                Value: &command.Parameter_S{
                        S: "fiona",
                },
        }

	nms := &command.Statement{
		Sql: "SELECT * FROM foo WHERE name=?",
		Value: []*command.Parameter{nmv},
	}

	nmqc := &command.NewQueryCommand{
		Statements: []*command.Statement{nms},
	}

	fmt.Println("Size of NewQueryCommand in JSON:", mustSizeofJSON(nmqc))
	fmt.Println("Size of NewQueryCommand in Proto:", mustSizeofNewQueryProto(nmqc))

        fmt.Println("===========================")
        fmt.Println("New model multi")
        fmt.Println("===========================")

        nmv = &command.Parameter{
                Value: &command.Parameter_S{
                        S: "fiona",
                },
        }

        nms = &command.Statement{
                Sql: "SELECT * FROM foo WHERE name=?",
                Value: []*command.Parameter{nmv},
        }

	nmss := make([]*command.Statement, 50)
	for i := 0; i < 50; i++ {
		nmss[i] = nms
	}

        nmqc = &command.NewQueryCommand{
                Statements: nmss,
        }

	// Encode statements as JSON and then compress.
        var nmsb bytes.Buffer
        nmgz, err := gzip.NewWriterLevel(&nmsb, gzip.BestCompression)
        if err != nil {
                log.Fatal(err)
        }
        if err := json.NewEncoder(nmgz).Encode(nmss); err != nil {
                log.Fatalf("failed to JSON encode and compress: %s", err.Error())
        }
        nmgz.Close()
	nmqz := &command.NewQueryCommand{
		CompressedStatements: nmsb.Bytes(),
        }

        fmt.Println("Size of NewQueryCommand in JSON:", mustSizeofJSON(nmqc))
        fmt.Println("Size of NewQueryCommand in Proto:", mustSizeofNewQueryProto(nmqc))
        fmt.Println("Size of NewQueryCommand in Proto, compressed:", mustSizeofNewQueryProto(nmqz))
}

func mustJSONMarshal(o interface{}) []byte {
	b, err := json.Marshal(o)
        if err != nil {
                panic("failed to marshal JSON")
        }
	return b
}

func mustSizeofJSON(o interface{}) int {
        b, err := json.Marshal(o)
        if err != nil {
                panic("failed to marshal JSON")
        }

	return len(b)
}

func mustSizeofNewQueryProto(c *command.NewQueryCommand) int {
	b, err := proto.Marshal(c)
        if err != nil {
		panic("failed to marshal protobuf")
        }
	return len(b)
}
