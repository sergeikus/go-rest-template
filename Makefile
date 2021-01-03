.PHONY: build test

project-name = go-server

help:
	@echo "'${project-name}' makefile"
	@echo "Available rules: "
	@echo "    build  - creates '${project-name}' image"
	@echo "    clean  - clean server"
	@echo "    up  - puts server container up"
	@echo "    down  - puts server container down"
	@echo "    start  - starts server container"
	@echo "    stop  - stops server container"
	@echo "    restart  - restarts server container"
	@echo

build: 
	@echo "=====   Building '${project-name}' image   ====="
	docker build -f build/Dockerfile . -t ${project-name} --force-rm
	- docker rmi -f `docker images -f "dangling=true" -q`
	@echo "=====   Image is created   ====="
	@echo

binary-locally:
	@echo "=====   Building '${project-name}' binary   ====="
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux go build -o bin/${project-name} cmd/main.go 
	@echo "=====   Successfully build '${project-name}'   ====="
	@echo	

up:
	@echo "==========   Starting '${project-name}' container   =========="
	docker run -it -d --network=host --name ${project-name} ${project-name}
	@echo "==========   '${project-name}' is UP   =========="
	@echo

stop:
	@echo "==========   Stopping '${project-name}' container   =========="
	docker stop ${project-name}
	@echo "==========   '${project-name}' is STOPPED   =========="
	@echo

start:
	@echo "==========   Starting '${project-name}' container   =========="
	docker start ${project-name}
	@echo "==========   Started '${project-name}'   =========="
	@echo

restart:
	@echo "==========   Restarting '${project-name}' container   =========="
	docker restart ${project-name}
	@echo "==========   Restarted '${project-name}'   =========="
	@echo

down:
	@echo "==========   Putting down '${project-name}' container   =========="
	- docker rm -fv ${project-name}
	@echo "==========   Put down '${project-name}'   =========="
	@echo

test:
	@echo "==========   Performing '${project-name}' unit testing with coverage report   =========="
	docker build -f build/Dockerfile-test . -t ${project-name}-test --force-rm
	- docker rmi -f `docker images -f "dangling=true" -q`
	docker create -ti --name ${project-name}-test ${project-name}-test bash
	docker cp ${project-name}-test:/server/coverage.html .
	docker rm -f ${project-name}-test
	@echo "==========   Go coverage report is available at 'server/coverage.html'   =========="


test-locally:
	@echo "==========   Performing '${project-name}' unit testing with coverage report   =========="
	@env CGO_ENABLED=0 go test -timeout 30s -coverprofile=c.out `go list ./... | grep -v cmd`
	@go tool cover -html=c.out -o coverage.html

clean:
	- $(MAKE) down
	- docker rmi -f ${project-name}
	rm -rf bin vendor
