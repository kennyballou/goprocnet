VERSION=0
MAJOR=0
MINOR=1
BUILD=`git rev-parse --short HEAD`
LDFLAGS="-X goprocnet.Version ${VERSION}.${MAJOR}.${MINOR}-${BUILD}"
all: build

build:
	@mkdir -p _build
	@go build -ldflags ${LDFLAGS} -o _build/goprocnet

test:
	@go test

install:
	@go install -ldflags ${LDFLAGS}

clean:
	@/bin/rm -rf _build/
	@/bin/rm -rf _dist/
