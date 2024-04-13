all: build 

build:
	go build -o git-sync .

test:
	go test -v ./...

clean:
	rm -f git-sync 
