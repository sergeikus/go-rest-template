.PHONY: test
SHELL:=/bin/bash

server-project-name = go-server
database-project-name = database

help:
	@echo "Available rules: "
	@echo "    Server rules:"
	@echo "        server-build  - creates '${server-project-name}' image"
	@echo "        server-clean  - clean server related stuff"
	@echo "        server-up  - sets server container up "
	@echo "                     NB: it's possible to overwrite database type defined"
	@echo "                         in the 'config.yaml'."
	@echo "                         Just provide: env DB_TYPE_INMEMORY=true before 'server-up'"
	@echo "                         This variable will force server to create an 'in-memory' database."
	@echo "                     IMPORTANT: by default server will try to connect to an external database,"
	@echo "                                so build and run database container before the server is up"
	@echo "        server-down  - puts server container down"
	@echo "        server-start  - starts server container"
	@echo "        server-stop  - stops server container"
	@echo "        server-restart  - restarts server container"
	@echo
	@echo "    Database rules:"
	@echo "        database-build  - creates database image"
	@echo "        database-clean  - clean database related stuff"
	@echo "        database-up  - sets database container up"
	@echo "        database-down  - puts database container down"
	@echo "        database-start  - starts database container"
	@echo "        database-stop  - stops database container"
	@echo "        database-restart  - restarts database container"
	@echo 
	@echo "    Test rules:"
	@echo "        test - performs Go unit testing with coverage report inside a docker container"
	@echo "        test-locally - performs Go unit testing with coverage report with host machine Go"
	@echo 
	@echo "    Other:"
	@echo "        clean - performs repository clean up, removes containers, images and temporary files"
	@echo

server-build: 
	@echo "=====   Building '${server-project-name}' image   ====="
	docker build -f build/Dockerfile . -t ${server-project-name}
	- docker rmi -f `docker images -f "dangling=true" -q`
	@echo "=====   Image is created   ====="
	@echo

server-binary-locally:
	@echo "=====   Building '${server-project-name}' binary   ====="
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux go build -o bin/${server-project-name} cmd/main.go 
	@echo "=====   Successfully build '${server-project-name}'   ====="
	@echo	

server-up:
	@echo "==========   Starting '${server-project-name}' container   =========="
	docker run --env DB_TYPE_INMEMORY=$(shell echo $$DB_TYPE_INMEMORY) -it -d --network=host --name ${server-project-name} ${server-project-name}
	@echo "==========   '${server-project-name}' is UP   =========="
	@echo

server-stop:
	@echo "==========   Stopping '${server-project-name}' container   =========="
	docker stop ${server-project-name}
	@echo "==========   '${server-project-name}' is STOPPED   =========="
	@echo

server-start:
	@echo "==========   Starting '${server-project-name}' container   =========="
	docker start ${server-project-name}
	@echo "==========   Started '${server-project-name}'   =========="
	@echo

server-restart:
	@echo "==========   Restarting '${server-project-name}' container   =========="
	docker restart ${server-project-name}
	@echo "==========   Restarted '${server-project-name}'   =========="
	@echo

server-down:
	@echo "==========   Putting down '${server-project-name}' container   =========="
	- docker rm -fv ${server-project-name}
	@echo "==========   Put down '${server-project-name}'   =========="
	@echo

server-clean:
	- $(MAKE) server-down
	- docker rmi -f ${server-project-name}
	rm -rf bin vendor c.out coverage.html

test:
	@echo "==========   Performing '${server-project-name}' unit testing with coverage report   =========="
	docker build -f build/Dockerfile-test . -t ${server-project-name}-test
	- docker rmi -f `docker images -f "dangling=true" -q`
	docker create -ti --name ${server-project-name}-test ${server-project-name}-test bash
	docker cp ${server-project-name}-test:/server/coverage.html .
	docker rm -f ${server-project-name}-test
	@echo "==========   Go coverage report is available in 'coverage.html'   =========="

test-locally:
	@echo "==========   Performing '${server-project-name}' unit testing with coverage report   =========="
	@env CGO_ENABLED=0 go test -timeout 30s -coverprofile=c.out `go list ./... | grep -v cmd`
	@go tool cover -html=c.out -o coverage.html
	@echo "==========   Go coverage report is available in 'coverage.html'   =========="

database-build:
	@echo "=====   Building '${server-project-name}' image   ====="
	docker build -f build/database/Dockerfile-postgres . -t ${database-project-name}
	- docker rmi -f `docker images -f "dangling=true" -q`
	@echo "=====   Image is created   ====="
	@echo

database-down:
	@echo "==========   Putting down '${database-project-name}' container   =========="
	- docker rm -fv ${database-project-name}
	@echo "==========   Put down '${database-project-name}'   =========="
	@echo

database-up:
	@echo "==========   Starting '${database-project-name}' container   =========="
	docker run -it -d --network=host --name ${database-project-name} ${database-project-name}
	@echo "==========   '${database-project-name}' is UP   =========="
	@echo

database-clean:
	- $(MAKE) database-down
	- docker rmi -f ${database-project-name}

clean: database-clean server-clean
	- docker rmi -f go-server-test
