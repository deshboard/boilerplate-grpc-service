PROTO_PATH = vendor/github.com/deshboard/boilerplate-proto/proto

.PHONY: proto

proto: ## Generate code from protocol buffer
	@mkdir -p apis
	protoc -I ${PROTO_PATH} ${PROTO_PATH}/boilerplate/boilerplate.proto --go_out=plugins=grpc:apis
	protoc -I ${PROTO_PATH} ${PROTO_PATH}/boilerplate2/boilerplate2.proto --go_out=plugins=grpc:apis

envcheck::
	$(call executable_check,protoc,protoc)
	$(call executable_check,protoc-gen-go,protoc-gen-go)
