# generate API request/response models based on proto file description
generate:
	protoc --proto_path=. --go_out=pkg/user-management-api --go_opt=paths=source_relative user-management-api.proto

dev-up:
	cd ./scripts/dev; \
	docker-compose up -d

dev-down:
	cd ./scripts/dev; \
	docker-compose down