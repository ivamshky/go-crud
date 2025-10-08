PROTO_PATH := grpc

GEN_GO_DIR := gen/grpc

.PHONY: proto clean

proto: $(shell find $(PROTO_PATH) -name "*.proto")
	@echo "-> Generating Go code from Protobufs..."

	protoc \
		--proto_path=$(PROTO_PATH) \
		--go_out=$(GEN_GO_DIR) \
		--go_opt=paths="source_relative" \
		--go-grpc_out=$(GEN_GO_DIR) \
		--go-grpc_opt=paths="source_relative" \
		$(shell find $(PROTO_PATH) -name "*.proto")

	@echo "-> Generation complete in $(GEN_GO_DIR)"

clean:
	@echo "-> Cleaning generated files..."
	rm -rf $(GEN_GO_DIR)/*
	@echo "-> Clean complete."