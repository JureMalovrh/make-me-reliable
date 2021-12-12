.PHONY: dev prod full remove test lint ci deps
dev: deps
	docker build . --target=build-env -t reliable-api:dev

prod: dev
	docker build . -t reliable-api:latest

full: dev prod

clean:
	docker image rm reliable-api:dev
	docker image rm reliable-api:latest

test-deps:
	docker-compose up -d mongo

test: dev test-deps
	docker-compose run --rm reliable-api-dev go test -v -race ./...

lint: dev
	docker run reliable-api:dev golangci-lint run ./...

ci: test lint

deps:
	docker build ./unreliable-api -f ./unreliable-api/Dockerfile -t unreliable-api:latest