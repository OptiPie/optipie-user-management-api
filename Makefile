generate:
	protoc --proto_path=. --go_out=pkg/user-management-api --go_opt=paths=source_relative user-management-api.proto