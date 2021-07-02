# snarkOSBlockTest

To run clone this repo, cd to it and then execute:

```
go get ./...
go run blocktest.go
```

In order to test with more than 1024 requests per second, you may need to increase the file descriptor limit via something like:

```
ulimit -n 20000
```
