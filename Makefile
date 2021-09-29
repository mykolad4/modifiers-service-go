build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/handler handler/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock