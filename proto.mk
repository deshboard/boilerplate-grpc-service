.PHONY: proto

proto: ## Generate code from protocol buffer
	@mkdir -p model
	protowrap -I ${PROTO_PATH} ${PROTO_PATH}/boilerplate/boilerplate.proto ${PROTO_PATH}/boilerplate2/boilerplate2.proto  --go_out=plugins=grpc:model

envcheck::
	$(call executable_check,protoc,protoc)
	$(call executable_check,protowrap,protowrap)
