install:
	mkdir /etc/leprechaun
	mkdir /var/log/leprechaun/
	mkdir /var/run/leprechaun/
	touch /var/log/leprechaun/info.log
	touch /var/log/leprechaun/error.log
	touch /var/run/leprechaun/.pid
	go build

uninstall:
	rm -rf /etc/leprechaun
	rm -rf /var/log/leprechaun
	rm -rf /var/run/leprechaun

build:
	cd bin/ && go build

format:
	gofmt -s -w .

test:
	go vet
	cd client && go vet
	cd log && go vet
	cd recipe && go vet
	cd recipe/schedule && go vet
	go test