PROTO_PATH = vendor/github.com/deshboard/boilerplate-proto/proto

.PHONY: proto

proto: ## Generate code from protocol buffer
	@mkdir -p proto
	protoc -I ${PROTO_PATH} ${PROTO_PATH}/boilerplate/boilerplate.proto --go_out=plugins=grpc:proto
	protoc -I ${PROTO_PATH} ${PROTO_PATH}/boilerplate2/boilerplate2.proto --go_out=plugins=grpc:proto

envcheck::
	$(call executable_check,protoc,protoc)
	$(call executable_check,protoc-gen-go,protoc-gen-go)
