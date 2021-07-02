package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var url = "http://128.199.15.82:3030/"
var concurrentRequests = 1000
var blockHash = "9b395f3efb734c5c4a7f0145be92a0f717cc37e24d4185997de0bde5a27e44f4"

func main() {
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
			return
		}
	}

}

func do(i int) error {
	if i%1000 == 0 {
		fmt.Printf("%v: %v\n", i, time.Now())
	}
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
