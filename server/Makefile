postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:alpine

redis:
	docker run --name redis -d -p 6379:6379 redis:alpine redis-server --save 60 1

createdb:
	docker exec -it postgres createdb --username=root --owner=root mirage

dropdb:
	docker exec -it postgres dropdb mirage

recreate:
	make dropdb && make createdb

start:
	docker start postgres && docker start redis

test:
	go test -v -cover ./service/... ./handler/...

e2e:
	go test -run '^\QTestMain_E2E\E$\'

lint:
	golangci-lint run

build:
	go build github.com/sentrionic/mirage

fmt:
	go fmt github.com/sentrionic/...