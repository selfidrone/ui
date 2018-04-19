VERSION=0.5
NAMESPACE=nicholasjackson

mocks:
	go generate ./...

test:
	GOMAXPROCS=7 go test -parallel 7 -cover -race ./...

dev:
	buffalo dev

build:
	go build -o drone-live .

go_relaser:
	goreleaser --snapshot --rm-dist --skip-validate

build_linux:
	CGO_ENABLED=0 GOOS=linux go build -o drone-live .

build_docker: build_linux
	docker build -t ${NAMESPACE}/drone-live:${VERSION} .
	docker tag ${NAMESPACE}/drone-live:${VERSION} ${NAMESPACE}/drone-live:${VERSION}
	docker tag ${NAMESPACE}/drone-live:${VERSION} ${NAMESPACE}/drone-live:latest

push_docker: 
	docker push ${NAMESPACE}/drone-live:${VERSION}
	docker push ${NAMESPACE}/drone-live:latest

run_docker:
	docker run -it -p 4000:4000 ${NAMESPACE}/drone-live:latest -nats nats://nats.drone.demo.gs:4222 -source /home/dronelive/

all_arch:
	GOOS=linux buffalo build -o bin/web_ui
	GOOS=linux GOARCH=arm GOARM=6 buffalo build -o bin_arm/web_ui
