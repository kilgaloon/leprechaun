export GO111MODULE=on
DESTDIR?=/

install:
	mkdir -p $(DESTDIR)/etc/leprechaun/{configs,recipes}
	sed -e "s#@@USER_HOME@@#$(HOME)#g;s#@@LEPRECHAUN_HOME@@#$(DESTDIR)#g" dist/configs/config.ini > $(DESTDIR)/etc/leprechaun/configs/config.ini
	cp dist/configs/debug_config.ini $(DESTDIR)/etc/leprechaun/configs/
	mkdir -p $(DESTDIR)/var/log/leprechaun/{server,workers.output}
	touch $(DESTDIR)/var/log/leprechaun/info.log
	touch $(DESTDIR)/var/log/leprechaun/error.log
	cp leprechaun.service /etc/systemd/system/
	go install ./cmd/leprechaun
	go install ./cmd/lepretools

install-user: install
	# This part depends on who install it
	# i assume it is user who want to run it
	mkdir -p ~/.config/systemd/user/
	sed -e "s#@@USER_HOME@@#$(HOME)#g;s#@@LEPRECHAUN_HOME@@#$(DESTDIR)#g" leprechaun-user.service > ~/.config/systemd/user/leprechaun-user.service

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
	cd config && go vet
	cd context && go vet
	cd event && go vet
	cd workers && go vet
	cd server && go vet
	cd agent && go vet
	cd api && go vet
	go test ./client -coverprofile=./client/coverage.txt -covermode=atomic
	go test ./config -coverprofile=./config/coverage.txt -covermode=atomic
	go test ./context -coverprofile=./context/coverage.txt -covermode=atomic
	go test ./event -coverprofile=./event/coverage.txt -covermode=atomic
	go test ./log -coverprofile=./log/coverage.txt -covermode=atomic
	go test ./workers -coverprofile=./workers/coverage.txt -covermode=atomic
	go test ./server -coverprofile=./server/coverage.txt -covermode=atomic
	go test ./api -coverprofile=./api/coverage.txt -covermode=atomic
	go test ./agent -coverprofile=./agent/coverage.txt -covermode=atomic
	go test ./recipe -coverprofile=./api/coverage.txt -covermode=atomic
	go test ./cron -coverprofile=./cron/coverage.txt -covermode=atomic
	go test ./notifier -coverprofile=./notifier/coverage.txt -covermode=atomic

test-with-report:
	go vet ./cmd/leprechaun
	go vet ./cmd/lepretools
	cd client && go vet
	cd log && go vet
	cd recipe && go vet
	cd recipe/schedule && go vet
	mkdir -p coverprofile
	go test ./client -coverprofile coverprofile/client.out
	go tool cover -html=coverprofile/client.out -o coverprofile/client.html
	go test ./config -coverprofile coverprofile/config.out
	go tool cover -html=coverprofile/config.out -o coverprofile/config.html
	go test ./context -coverprofile coverprofile/context.out
	go tool cover -html=coverprofile/context.out -o coverprofile/context.html
	go test ./event -coverprofile coverprofile/event.out
	go tool cover -html=coverprofile/event.out -o coverprofile/event.html
	go test ./log -coverprofile coverprofile/log.out
	go tool cover -html=coverprofile/log.out -o coverprofile/log.html
	go test ./workers -coverprofile coverprofile/workers.out
	go tool cover -html=coverprofile/workers.out -o coverprofile/workers.html
	go test ./server -coverprofile coverprofile/server.out
	go tool cover -html=coverprofile/server.out -o coverprofile/server.html
	go test ./api -coverprofile coverprofile/api.out
	go tool cover -html=coverprofile/api.out -o coverprofile/api.html
	go test ./agent -coverprofile coverprofile/agent.out
	go tool cover -html=coverprofile/agent.out -o coverprofile/agent.html
	go test ./recipe -coverprofile coverprofile/recipe.out
	go tool cover -html=coverprofile/recipe.out -o coverprofile/recipe.html
	go test ./cron -coverprofile coverprofile/cron.out
	go tool cover -html=coverprofile/cron.out -o coverprofile/cron.html
	go test ./notifier -coverprofile coverprofile/notifier.out
	go tool cover -html=coverprofile/notifier.out -o coverprofile/notifier.html
	go test ./notifier/notifications -coverprofile coverprofile/notifications.out
	go tool cover -html=coverprofile/notifications.out -o coverprofile/notifications.html