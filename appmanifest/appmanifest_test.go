package main

import (
	"bytes"
	"testing"
)

func TestCalucalteMD5s(t *testing.T) {
	ok := []string{
		"f1c9645dbc14efddc7d8a322685f26eb",
		"f1c9645dbc14efddc7d8a322685f26eb",
		"f1c9645dbc14efddc7d8a322685f26eb",
		"93b885adfe0da089cdf634904fd59f71",
	}
	buf := make([]byte, MaxChunkSize*3)
	r := bytes.NewReader(buf)
	md5s, err := calculateMD5s(r)
	if err != nil {
		t.Fatal(err)
	}

	for i, h := range md5s {
		if ok[i] != h {
			t.Fatal("expected", ok[i], "got", h)
		}
	}
}

/*
Benchmark10MB-8  	    2000	    813545 ns/op	10485878 B/op	      12 allocs/op
Benchmark100MB-8 	    2000	    858545 ns/op	10485772 B/op	       2 allocs/op
Benchmark1000MB-8	       1	1513403113 ns/op	10503832 B/op	     623 allocs/op

Benchmark10MB-8  	  500000	      3297 ns/op	   32896 B/op	       3 allocs/op
Benchmark100MB-8 	  300000	      3528 ns/op	   32897 B/op	       3 allocs/op
Benchmark1000MB-8	       1	1967685209 ns/op	 3334064 B/op	     829 allocs/op
*/

func benchmarkSize(b *testing.B, size int) {
	var buf = make([]byte, size)
	r := bytes.NewReader(buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calculateMD5s(r)
	}
}

func Benchmark10MB(b *testing.B) {
	benchmarkSize(b, MaxChunkSize)
}

func Benchmark100MB(b *testing.B) {
	benchmarkSize(b, MaxChunkSize*10)
}

func Benchmark1000MB(b *testing.B) {
	benchmarkSize(b, MaxChunkSize*100)
}
