BINARY_NAME=drive-scanner
VERSION=${DS_VER}
DATE=$(shell date -u +'%Y-%m-%d %I:%M:%S%p %Z')

build:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME} -ldflags "-X 'main.version=${VERSION}' -X 'main.date=${DATE}'" main.go

run:
	./${BINARY_NAME}

build_and_run: build run

clean:
	go clean
	rm ${BINARY_NAME}
