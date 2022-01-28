build:
	go build -o bin/MeasMon src/main.go

run:
	go run src/main.go


compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin/main-linux-arm src/main.go
	GOOS=linux GOARCH=arm64 go build -o bin/main-linux-arm64 src/main.go
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 src/main.go

all: hello build
