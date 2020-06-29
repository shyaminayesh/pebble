all: run


run:
	rm -rf pebble
	go build ./pebble.go
	./pebble schema migrate