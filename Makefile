PREFIX  ?= /usr/local

qcal: clean
	GOOS=linux GOARCH=amd64 go build -o qcard -ldflags="-s -w"

linux-arm:
	GOOS=linux GOARCH=arm go build -o qcard -ldflags="-s -w"

darwin:	
	GOOS=darwin GOARCH=amd64 go build -o qcard -ldflags="-s -w"

windows:
	GOOS=windows GOARCH=amd64 go build -o qcard.exe -ldflags="-s -w"

clean:
	rm -f qcard

install: 
	install -d $(PREFIX)/bin/
	install -m 755 qcard $(PREFIX)/bin/qcard

uninstall:
	rm -f $(PREFIX)/bin/qcard
