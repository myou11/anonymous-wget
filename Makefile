all: build

build: ss awget
	echo "build complete"

ss:
	go build -o ss ./cmd/ss/ss.go

awget:
	go build -o awget ./cmd/awget/awget.go

clean:
	rm -f awget ss

test: build
    #this is localhost test
	echo "firing up 3 ss instances at ports 8989 8990 and 8991"
	./ss -p 8989 &
	./ss -p 8990 &
	./ss -p 8991 &
	echo "issuing a request to https://www.torproject.org/index.html"
	./awget -c examples/chaingang.txt https://www.torproject.org/index.html
	echo "terminating ss instances"
	pkill ss
