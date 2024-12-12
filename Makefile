.PHONY: build run test-start test-stop test-ci

default: build

build:
	docker buildx build --platform linux/amd64 --format docker -t vm75/easy-share .

run:
	docker run --rm -p 80:80 -p 137:137 -p 138:138 -p 139:139 -p 445:445 -p 2049:2049 vm75/easy-share

test-start:
	./test/cmd.sh run

test-stop:
	./test/cmd.sh stop

test-ci:
	act -s DOCKER_USERNAME -s DOCKER_PASSWORD -s GITHUB_TOKEN
