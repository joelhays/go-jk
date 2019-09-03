all: clean test build
clean:
		go clean ./...
		rm -f ./go-jk.exe
test:
		go test -v ./...
build:
		go build -v
run: build
		./go-jk.exe $(ARGS)
pprof:
		go tool pprof -http=:8080 go-jk.prof