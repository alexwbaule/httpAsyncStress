
PACKAGE:=stress
OUTPUT:=stress

.PHONY: linux windows mac build

linux: GOOS=linux
linux: GOARCH=amd64
linux: build

windows: GOOS=windows
windows: GOARCH=amd64
windows: OUTPUT=stress.exe
windows: build

mac: GOOS=darwin
mac: GOARCH=amd64
mac: build

build:
	mkdir -p bin/${GOOS}/${GOARCH}
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/${GOOS}/${GOARCH}/$(OUTPUT) $(PACKAGE).go
	cp -f bin/${GOOS}/${GOARCH}/$(OUTPUT) $(PACKAGE)
	chmod 755 $(PACKAGE)
