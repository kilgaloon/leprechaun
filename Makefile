install:
	mkdir /etc/leprechaun
	mkdir /etc/leprechaun/recipes
	cp -r configs /etc/leprechaun/configs
	mkdir /var/log/leprechaun/
	mkdir /var/log/leprechaun/server
	mkdir /var/run/leprechaun/
	touch /var/log/leprechaun/info.log
	touch /var/log/leprechaun/server/info.log
	touch /var/log/leprechaun/error.log
	touch /var/log/leprechaun/server/error.log
	touch /var/run/leprechaun/.pid
	go install

uninstall:
	rm -rf /etc/leprechaun
	rm -rf /var/log/leprechaun
	rm -rf /var/run/leprechaun

build:
	go build

rebuild:
	go clean
	go build

format:
	gofmt -s -w .

test:
	go vet
	cd client && go vet
	cd log && go vet
	cd recipe && go vet
	cd recipe/schedule && go vet
	go test