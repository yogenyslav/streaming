include .env

.PHONY: docker_build
docker_build:
	cd $(src) && docker build -t ${PROJECT_DIR}-$(image) .
	docker tag ${PROJECT_DIR}-$(image) $(shell whoami)/${PROJECT_DIR}-$(image):$(tag)
	docker push $(shell whoami)/${PROJECT_DIR}-$(image):$(tag)

.PHONY: docker_up
docker_up:
	docker compose up -d --build

.PHONY: docker_down
docker_down:
	docker compose down

.PHONY: docker_remove
docker_remove: docker_down
	docker volume rm ${PROJECT_DIR}_pg_data
	docker volume rm ${PROJECT_DIR}_mongo
	docker volume rm ${PROJECT_DIR}_mongo_conf
	docker volume rm ${PROJECT_DIR}_kafka_conf
	docker volume rm ${PROJECT_DIR}_kafka_data
	docker volume rm ${PROJECT_DIR}_kafka_secrets
	docker volume rm ${PROJECT_DIR}_static
	docker volume rm ${PROJECT_DIR}_zoo_data
	docker volume rm ${PROJECT_DIR}_minio
	docker image rm ${PROJECT_DIR}_api
	docker image rm ${PROJECT_DIR}_frame_service
	docker image rm ${PROJECT_DIR}_detection_service

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

.PHONY: proto_py
proto_py:
	python3 -m grpc_tools.protoc -Iprotos --python_out=frame_service/pb --pyi_out=frame_service/pb --grpc_python_out=frame_service/pb protos/frame.proto

.PHONY: proto_go
proto_go:
	protoc --go_out=api --go_opt=Mprotos/frame.proto=internal/pb \
		--go-grpc_out=api --go-grpc_opt=Mprotos/frame.proto=internal/pb \
		protos/frame.proto

.PHONY: kuber_deploy
kuber_deploy:
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/minio.yaml --namespace=streaming
	kubectl apply -f k8s/postgres.yaml --namespace=streaming
	kubectl apply -f k8s/mongo.yaml --namespace=streaming
	kubectl apply -f k8s/zookeeper.yaml --namespace=streaming
	kubectl apply -f k8s/kafka.yaml --namespace=streaming
	kubectl apply -f k8s/detection.yaml --namespace=streaming
	kubectl apply -f k8s/frame.yaml --namespace=streaming
	kubectl apply -f k8s/api.yaml --namespace=streaming
