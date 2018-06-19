install:
	mkdir /etc/leprechaun
	mkdir /var/log/leprechaun/
	mkdir /var/run/leprechaun/
	touch /var/log/leprechaun/info-client.log
	touch /var/log/leprechaun/error-client.log
	cd bin/ && go build -u leprechaun

uninstall:
	rm -rf /etc/leprechaun
	rm -rf /var/log/leprechaun
	rm -rf /var/run/leprechaun

build:
	cd bin/ && go build -o leprechaun

format:
	gofmt -s -w src/

test:
	cd bin && go vet
	cd src/client && go vet
	cd src/cmd && go vet
	cd src/log && go vet
	cd src/recipe && go vet
	cd src/recipe/schedule && go vet