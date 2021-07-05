package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var url = "http://128.199.15.82:3030/"
var concurrentRequests = 500
var blockHashes = []string{
	"1238cd972e89c3b687cf16b3aaaafb91f12f107f2554f861b21b6a51090b8e93",
	"ba16c5647ef81f19e3a7eb039ca71488c5fa3949d0d83e3e6f674de1a2250c34",
	"00521430d5d54bdc3b30657eb6b82858180f90878b302d3dc8f1cbd3271a1a98",
	"a4f675c4bbef63acf7c0b50b4bd7c0d34322e0f02fe8d5c43d0570df64f1f428",
	"4a0a43653aaec9384849c0775e0d24d3182e9d88a339b7dd0f982d712c9d7961",
}

func main() {
	rand.Seed(time.Now().Unix())
	fmt.Printf("blocktest\n")

	do(-1)

	for i := 0; ; i++ {
		g := errgroup.Group{}

		for j := 0; j < concurrentRequests; j++ {
			i := i
			j := j
			g.Go(func() error {
				return do(i*concurrentRequests + j)
			})
		}

		err := g.Wait()
		if err != nil {
			fmt.Printf("%+v", err)
		}
	}

}

func randomItem(slice []string) string {
	return slice[rand.Intn(len(slice))]
}

func do(i int) error {
	if i%1000 == 0 {
		fmt.Printf("%v: %v\n", i, time.Now())
	}

	blockHash := randomItem(blockHashes)
	res, err := http.Post(url, "application/json", strings.NewReader(`{"jsonrpc": "2.0","id":"documentation","method": "getblock","params": ["`+blockHash+`"]}`))
	if err != nil {
		return errors.Wrap(err, "")
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "")
	}

	responseString := string(bytes)

	if i == -1 {
		fmt.Printf("normal data\n")
		fmt.Printf("%v\n", responseString)
		return io.EOF
	}

	if strings.Contains(responseString, "null") || strings.Contains(responseString, "error") || strings.Contains(responseString, `"confirmations": 0`) {
		fmt.Printf("strange data: \n")
		fmt.Printf("%v: %v\n", i, time.Now())
		fmt.Printf("%v\n", responseString)
		return io.EOF
	}

	return nil
}
