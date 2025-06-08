package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "f", "file to read from")
	flag.StringVar(&to, "to", "t", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	fileStat, err := os.Stat("testdata/input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	if limit == 0 {
		limit = fileStat.Size() - offset
	}
	err = Copy(from, to, offset, limit)
	if err != nil {
		fmt.Println(err)
		return
	}
}
