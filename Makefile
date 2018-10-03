install:
	mkdir /etc/leprechaun
	mkdir /etc/leprechaun/recipes
	cp -r dist/configs /etc/leprechaun/configs
	mkdir /var/log/leprechaun/
	mkdir /var/log/leprechaun/server
	mkdir /var/run/leprechaun/
	touch /var/log/leprechaun/info.log
	touch /var/log/leprechaun/server/info.log
	touch /var/log/leprechaun/error.log
	touch /var/log/leprechaun/server/error.log
	touch /var/run/leprechaun/.pid
	go install ./cmd/leprechaun
	go install ./cmd/lepretools

uninstall:
	rm -rf /etc/leprechaun
	rm -rf /var/log/leprechaun
	rm -rf /var/run/leprechaun

build:
	go build ./cmd/leprechaun
	go build ./cmd/lepretools

rebuild:
	go clean
	make build

format:
	gofmt -s -w .

test:
	go vet ./cmd/leprechaun
	go vet ./cmd/lepretools
	cd client && go vet
	cd log && go vet
	cd recipe && go vet
	cd recipe/schedule && go vet
	go test ./client -cover
	go test ./config -cover
	go test ./context -cover
	go test ./event -cover
	go test ./log -cover
	go test ./workers -cover
	go test ./server -cover
	go test ./api -cover

test-with-report:
	go vet ./cmd/leprechaun
	go vet ./cmd/lepretools
	cd client && go vet
	cd log && go vet
	cd recipe && go vet
	cd recipe/schedule && go vet
	mkdir -p coverprofile
	go test client -coverprofile coverprofile/client.out
	go tool cover -html=coverprofile/client.out -o coverprofile/client.html
	go test config -coverprofile coverprofile/config.out
	go tool cover -html=coverprofile/config.out -o coverprofile/config.html
	go test context -coverprofile coverprofile/context.out
	go tool cover -html=coverprofile/context.out -o coverprofile/context.html
	go test event -coverprofile coverprofile/event.out
	go tool cover -html=coverprofile/event.out -o coverprofile/event.html
	go test log -coverprofile coverprofile/log.out
	go tool cover -html=coverprofile/log.out -o coverprofile/log.html
	go test workers -coverprofile coverprofile/workers.out
	go tool cover -html=coverprofile/workers.out -o coverprofile/workers.html
	go test server -coverprofile coverprofile/server.out
	go tool cover -html=coverprofile/server.out -o coverprofile/server.html
	go test api -coverprofile coverprofile/api.out
	go tool cover -html=coverprofile/api.out -o coverprofile/api.html
	go test agent -coverprofile coverprofile/agent.out
	go tool cover -html=coverprofile/agent.out -o coverprofile/agent.html