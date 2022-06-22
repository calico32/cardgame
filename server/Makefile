LD_FLAGS = -ldflags "\
	-X build.mode=release \
	-X build.commit=$(shell git rev-parse HEAD) \
	-X build.branch=$(shell git rev-parse --abbrev-ref HEAD) \
	-X build.version=$(shell git describe --tags --always) \
	-X build.buildTime=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')"

cardgame-server: *.go
	go build $(LD_FLAGS) -o cardgame-server
