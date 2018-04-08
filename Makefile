install:
	mkdir /etc/rainbow
	mkdir /var/log/rainbow/
	mkdir /var/run/rainbow/
	go build

uninstall:
	rm -rf /etc/rainbow
	rm -rf /var/log/rainbow
	rm -rf /var/run/rainbow