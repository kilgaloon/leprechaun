install:
	mkdir /etc/leprechaun
	mkdir /var/log/leprechaun/
	mkdir /var/run/leprechaun/
	touch /var/log/leprechaun
	go build

uninstall:
	rm -rf /etc/leprechaun
	rm -rf /var/log/leprechaun
	rm -rf /var/run/leprechaun

format:
	gofmt -s -w src/

test:
	go vet