PROJECT?=opscore-workshop-admin
APP?=request-generator
PORT?=80
PORT_APP?=7784

CONTAINER_IMAGE?=$(PROJECT)/${APP}
RELEASE?=0.0.1

clean:
	rm -f bin/${APP}

gin:
	go get github.com/codegangsta/gin

dep:
	dep ensure

dep-tutu:
	dep ensure -update stash.tutu.ru/golang/readiness
	dep ensure -update stash.tutu.ru/golang/resources
	dep ensure -update stash.tutu.ru/golang/context_os
	dep ensure -update stash.tutu.ru/golang/http-server
	dep ensure -update stash.tutu.ru/golang/log
	dep ensure -update stash.tutu.ru/golang/envs
	dep ensure -update stash.tutu.ru/golang/opentracing

init: gin
	git config --global url."ssh://git@depot.tutu.ru:7999/".insteadOf "https://stash.tutu.ru/scm/"
	dep ensure -update

gorun: dep clean
	go build -o bin/${APP} -tags "dev load_envs" ./cmd/ && bin/${APP}

watcher: dep gin
	gin --build cmd/ --logPrefix watcher --immediate --buildArgs "-tags 'dev load_envs'" run

container:
	docker build -t $(CONTAINER_IMAGE):$(RELEASE) .

run: container
	docker stop $(CONTAINER_IMAGE):$(RELEASE) || true && docker rm $(CONTAINER_IMAGE):$(RELEASE) || true
	docker run --name ${APP} -p ${PORT}:${PORT_APP} --rm \
		-e "PORT=${PORT}" \
		$(CONTAINER_IMAGE):$(RELEASE)

test:
	go test -v -race ./...
