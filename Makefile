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
	docker volume rm ${PROJECT_DIR}_mongo_data
	docker volume rm ${PROJECT_DIR}_kafka1_data
	docker volume rm ${PROJECT_DIR}_kafka2_data
	docker volume rm ${PROJECT_DIR}_kafka3_data
	docker volume rm ${PROJECT_DIR}_zoo_data
	docker volume rm ${PROJECT_DIR}_minio_data
	docker image rm ${PROJECT_DIR}_api
	docker image rm ${PROJECT_DIR}_framer
	docker image rm ${PROJECT_DIR}_detection
	#docker image rm ${PROJECT_DIR}_responser
	docker image rm ${PROJECT_DIR}_orchestrator

.PHONY: docker_restart
docker_restart: docker_down docker_up

.PHONY: docker_purge_restart
docker_purge_restart: docker_remove docker_up

.PHONY: local
local:
	docker compose up pg frame_mongo detection_mongo kafka -d --build
	go run cmd/rest/main.go

.PHONY: migrate_up
migrate_up:
	cd api/migrations && goose postgres "user=${POSTGRES_USER} \
		password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=disable \
		host=${POSTGRES_HOST} port=${POSTGRES_PORT}" up

.PHONY: migrate_down
migrate_down:
	cd api/migrations && goose postgres "user=${POSTGRES_USER} \
		password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=disable \
		host=localhost port=${POSTGRES_PORT}" down

.PHONY: migrate_new
migrate_new:
	cd api/migrations && goose postgres "user=${POSTGRES_USER} \
		password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=disable \
		host=localhost port=${POSTGRES_PORT}" create $(name) sql

.PHONY: proto_go
proto_go:
	protoc --go_out=$(dest) --go_opt=Mprotos/streaming.proto=internal/pb \
		--go-grpc_out=$(dest) --go-grpc_opt=Mprotos/streaming.proto=internal/pb \
		protos/streaming.proto
