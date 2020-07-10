all: run


install:
	rm -rf pebble
	go build ./pebble.go
	sudo cp pebble /usr/bin/pebble

uninstall:
	sudo rm -rf /usr/bin/pebble

run:
	rm -rf pebble
	go build ./pebble.go
	./pebble seed