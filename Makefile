
# This how we want to name the binary output
SOURCES := tmswebsite.go
BINARY := tmswebsite.exe
COVERAGEFILE := cover.out

# These are the values we want to pass for VERSION and BUILD strings
# git tag 1.0.1
# git commit -am "One more change after the tags"
GIT_VERSION := $(shell git describe --tags)
ifeq ($(OS),Windows_NT)
BUILD_DATE := $(shell cmd /C date /T) 
else
BUILD_DATE := $(shell date +"%d-%m-%y")
endif
# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X main.Version=${GIT_VERSION} -X main.Build=${BUILD_DATE}"

# Builds the project and runs the coverage
ifeq ($(OS),Windows_NT)
all:	build test install buildLinux
else
all:	buildLinux test install
endif
# Installs our project: copies binaries
install:
	go install ${LDFLAGS}

build:
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY}
	
buildLinux:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY}
	
test:
	go test

cover:
	go test -coverprofile ${COVERAGEFILE}
	go tool cover -html ${COVERAGEFILE}


# Cleans our project: deletes binaries
clean:  if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install test cover