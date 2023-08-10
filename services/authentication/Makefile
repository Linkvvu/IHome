.PHONY: proto
proto:
	protoc --proto_path=./proto --micro_out=./proto --go_out=./proto		\
		   --micro_opt=paths=source_relative --go_opt=paths=source_relative	\
		   proto/authService.proto