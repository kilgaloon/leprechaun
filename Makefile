export GO111MODULE=on

debug:
	go run cmd/leprechaun/main.go --pid=./var/run/leprechaun/.pid --ini=./dist/configs/config.ini --debug=true

setup:
	mkdir /etc/leprechaun
	mkdir /etc/leprechaun/recipes
	cp dist/configs/config.ini /etc/leprechaun
	mkdir /var/run/leprechaun/
	mkdir /var/log/leprechaun/
	mkdir /var/log/leprechaun/workers.output
	touch /var/log/leprechaun/info.log
	touch /var/log/leprechaun/error.log

install:
	make setup
	go mod vendor
	go build -o $(GOPATH)/bin/leprechaun ./cmd/leprechaun
	go build -o $(GOPATH)/bin/lepretools ./cmd/lepretools

install-remote-service:
	make setup
	go build -tags remote -o $(GOPATH)/bin/leprechaunrmt ./cmd/leprechaun

uninstall:
	rm -rf /etc/leprechaun
	rm -rf /var/log/leprechaun
	rm -rf /var/run/leprechaun

build:
	go build ./cmd/leprechaun
	go build ./cmd/lepretools

build-remote-service:
	go build -tags remote -o ./leprechaunrmt ./cmd/leprechaun

#can be used to test secure connection between remote and client
self-ca:
	openssl req -newkey rsa:2048 -nodes -keyout key.pem -x509 -days 365 -out certificate.pem

rebuild:
	go clean
	make build

format:
	gofmt -s -w .

test-package:
	RUN_MODE=test go test -race ./${package} -coverprofile coverprofile/${package}.out -v
	go tool cover -html=coverprofile/${package}.out -o coverprofile/${package}.html

test-verbose:
	go vet ./cmd/leprechaun
	go vet ./cmd/lepretools
	cd client && go vet
	cd log && go vet
	cd recipe && go vet
	cd recipe/schedule && go vet
	cd config && go vet
	cd context && go vet
	cd workers && go vet
	cd server && go vet
	cd agent && go vet
	cd api && go vet
	RUN_MODE=test go test -race ./client -coverprofile=./client/coverage.txt -covermode=atomic -v
	RUN_MODE=test go test -race ./config -coverprofile=./config/coverage.txt -covermode=atomic -v 
	RUN_MODE=test go test -race ./context -coverprofile=./context/coverage.txt -covermode=atomic -v
	RUN_MODE=test go test -race ./log -coverprofile=./log/coverage.txt -covermode=atomic -v
	RUN_MODE=test go test -race ./workers -coverprofile=./workers/coverage.txt -covermode=atomic -v
	RUN_MODE=test go test -race ./server -coverprofile=./server/coverage.txt -covermode=atomic -v
	RUN_MODE=test go test -race ./api -coverprofile=./api/coverage.txt -covermode=atomic -v
	RUN_MODE=test go test -race ./agent -coverprofile=./agent/coverage.txt -covermode=atomic -v
	RUN_MODE=test go test -race ./recipe -coverprofile=./api/coverage.txt -covermode=atomic -v
	RUN_MODE=test go test -race ./cron -coverprofile=./cron/coverage.txt -covermode=atomic -v
	RUN_MODE=test go test -race ./notifier -coverprofile=./notifier/coverage.txt -covermode=atomic -v
	RUN_MODE=test go test -race ./notifier/notifications -coverprofile=./notifier/notifications/coverage.txt -v
	RUN_MODE=test go test -race ./daemon -coverprofile=./daemon/coverage.txt -v
	RUN_MODE=test go test -race ./remote -coverprofile=./remote/coverage.txt -v

test:
	go vet ./cmd/leprechaun
	go vet ./cmd/lepretools
	cd client && go vet
	cd log && go vet
	cd recipe && go vet
	cd recipe/schedule && go vet
	cd config && go vet
	cd context && go vet
	cd workers && go vet
	cd server && go vet
	cd agent && go vet
	cd api && go vet
	RUN_MODE=test go test -race ./client -coverprofile=./client/coverage.txt -covermode=atomic
	RUN_MODE=test go test -race ./config -coverprofile=./config/coverage.txt -covermode=atomic
	RUN_MODE=test go test -race ./context -coverprofile=./context/coverage.txt -covermode=atomic
	RUN_MODE=test go test -race ./log -coverprofile=./log/coverage.txt -covermode=atomic
	RUN_MODE=test go test -race ./workers -coverprofile=./workers/coverage.txt -covermode=atomic
	RUN_MODE=test go test -race ./server -coverprofile=./server/coverage.txt -covermode=atomic
	RUN_MODE=test go test -race ./api -coverprofile=./api/coverage.txt -covermode=atomic
	RUN_MODE=test go test -race ./agent -coverprofile=./agent/coverage.txt -covermode=atomic
	RUN_MODE=test go test -race ./recipe -coverprofile=./api/coverage.txt -covermode=atomic
	RUN_MODE=test go test -race ./cron -coverprofile=./cron/coverage.txt -covermode=atomic
	RUN_MODE=test go test -race ./notifier -coverprofile=./notifier/coverage.txt -covermode=atomic
	RUN_MODE=test go test -race ./notifier/notifications -coverprofile=./notifier/notifications/coverage.txt
	RUN_MODE=test go test -race ./daemon -coverprofile=./daemon/coverage.txt 
	RUN_MODE=test go test -race ./remote -coverprofile=./remote/coverage.txt 

test-with-report:
	go vet ./cmd/leprechaun
	go vet ./cmd/lepretools
	cd client && go vet
	cd log && go vet
	cd recipe && go vet
	cd recipe/schedule && go vet
	mkdir -p coverprofile
	RUN_MODE=test go test -race ./client -coverprofile coverprofile/client.out
	go tool cover -html=coverprofile/client.out -o coverprofile/client.html
	RUN_MODE=test go test -race ./config -coverprofile coverprofile/config.out
	go tool cover -html=coverprofile/config.out -o coverprofile/config.html
	RUN_MODE=test go test -race ./context -coverprofile coverprofile/context.out
	go tool cover -html=coverprofile/context.out -o coverprofile/context.html
	RUN_MODE=test go test -race ./log -coverprofile coverprofile/log.out
	go tool cover -html=coverprofile/log.out -o coverprofile/log.html
	RUN_MODE=test go test -race ./workers -coverprofile coverprofile/workers.out
	go tool cover -html=coverprofile/workers.out -o coverprofile/workers.html
	RUN_MODE=test go test -race ./server -coverprofile coverprofile/server.out
	go tool cover -html=coverprofile/server.out -o coverprofile/server.html
	RUN_MODE=test go test -race ./api -coverprofile coverprofile/api.out
	go tool cover -html=coverprofile/api.out -o coverprofile/api.html
	RUN_MODE=test go test -race ./agent -coverprofile coverprofile/agent.out
	go tool cover -html=coverprofile/agent.out -o coverprofile/agent.html
	RUN_MODE=test go test -race ./recipe -coverprofile coverprofile/recipe.out
	go tool cover -html=coverprofile/recipe.out -o coverprofile/recipe.html
	RUN_MODE=test go test -race ./cron -coverprofile coverprofile/cron.out
	go tool cover -html=coverprofile/cron.out -o coverprofile/cron.html
	RUN_MODE=test go test -race ./notifier -coverprofile coverprofile/notifier.out
	go tool cover -html=coverprofile/notifier.out -o coverprofile/notifier.html
	RUN_MODE=test go test -race ./notifier/notifications -coverprofile coverprofile/notifications.out
	go tool cover -html=coverprofile/notifications.out -o coverprofile/notifications.html
	RUN_MODE=test go test -race ./daemon -coverprofile coverprofile/daemon.out
	go tool cover -html=coverprofile/daemon.out -o coverprofile/daemon.html
	RUN_MODE=test go test -race ./remote -coverprofile coverprofile/remote.out
	go tool cover -html=coverprofile/remote.out -o coverprofile/remote.html
