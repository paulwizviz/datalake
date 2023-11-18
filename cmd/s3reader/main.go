package main

import (
	"fmt"
	"log"
	"os"

	"github.com/paulwizviz/datalake/internal/blockutil"
)

func main() {
	keys, err := blockutil.ReadS3ListURL(blockutil.S3URL)
	if err != nil {
		log.Fatal(err)
	}
	for _, key := range keys {
		o, _ := blockutil.ReadObjectByKey(key)
		n := blockutil.GetBlockNumber(o)
		g, _ := os.Create(fmt.Sprintf("./testdata/%v.pb", key))
		defer g.Close()
		f, _ := os.Create(fmt.Sprintf("./testdata/%v.pb", n))
		defer f.Close()
		f.Write(o)
		g.Write(o)
	}
}
