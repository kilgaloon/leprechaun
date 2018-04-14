install:
	mkdir /etc/leprechaun
	mkdir /var/log/leprechaun/
	mkdir /var/run/leprechaun/
	touch /var/log/leprechaun/info-client.log
	touch /var/log/leprechaun/error-client.log
	go build

uninstall:
	rm -rf /etc/leprechaun
	rm -rf /var/log/leprechaun
	rm -rf /var/run/leprechaun

format:
	gofmt -s -w src/

test:
	go vet