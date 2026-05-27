.PHONY: install build clean release

install:
	go install .

build:
	go build -o ddx .

clean:
	rm -f ddx

release:
	./scripts/release.sh
