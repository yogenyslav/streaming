include .env

.PHONY: docker_up
docker_up:
	docker compose up -d --build

.PHONY: docker_down
docker_down:
	docker compose down

.PHONY: docker_remove
docker_remove: docker_down
	docker volume rm ${PROJECT_DIR}_pg_data
	docker volume rm ${PROJECT_DIR}_prom_data
	docker volume rm ${PROJECT_DIR}_jaeger_data
	docker volume rm ${PROJECT_DIR}_s3_data
	docker volume rm ${PROJECT_DIR}_zoo_data
	docker volume rm ${PROJECT_DIR}_kafka_data
	docker volume rm ${PROJECT_DIR}_mongo_data
	docker image rm api
	docker image rm orchestrator

.PHONY: docker_restart
docker_restart: docker_down docker_up

.PHONY: docker_purge_restart
docker_purge_restart: docker_remove docker_up

.PHONY: migrate_up
migrate_up:
	cd migrations && goose postgres "user=${POSTGRES_USER} \
		password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=disable \
		host=${POSTGRES_HOST} port=${POSTGRES_PORT}" up

.PHONY: migrate_down
migrate_down:
	cd migrations && goose postgres "user=${POSTGRES_USER} \
		password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=disable \
		host=localhost port=${POSTGRES_PORT}" down

.PHONY: migrate_new
migrate_new:
	cd migrations && goose create $(name) sql

.PHONY: proto_go
proto_go:
	protoc --proto_path=./proto --go_out=api \
        		--go-grpc_out=api proto/streaming/api.proto proto/streaming/enum.proto
	protoc --proto_path=./proto --go_out=orchestrator \
				--go-grpc_out=orchestrator proto/streaming/api.proto proto/streaming/enum.proto

.PHONY: tests
tests:
	docker compose -f test-docker-compose.yaml up pg -d
	sleep 5
	cd migrations && goose postgres "user=${POSTGRES_TEST_USER} \
    		password=${POSTGRES_TEST_PASSWORD} dbname=${POSTGRES_TEST_DB} sslmode=disable \
    		host=${POSTGRES_TEST_HOST} port=${POSTGRES_TEST_PORT}" up
	cd api && go test ./internal/streaming/query/... -cover -v -tags=integration
	cd api && go test ./internal/streaming/response/... -cover -v -tags=integration
	docker compose -f test-docker-compose.yaml down
	sleep 5
	docker volume rm ${PROJECT_DIR}_pg_test_data