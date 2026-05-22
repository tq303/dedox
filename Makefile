.PHONY: install build clean

install:
	go install .

build:
	go build -o ddx .

clean:
	rm -f ddx
