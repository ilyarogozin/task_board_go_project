# =========================
# CONFIG
# =========================

PROTO_DIR := proto
GEN_DIR   := gen/go

PROTO_FILES := $(shell find $(PROTO_DIR) -name "*.proto")

GO_OUT        := $(GEN_DIR)
GO_GRPC_OUT   := $(GEN_DIR)

# =========================
# TOOLS
# =========================

PROTOC := protoc
PROTOC_GEN_GO := protoc-gen-go
PROTOC_GEN_GO_GRPC := protoc-gen-go-grpc

# =========================
# TARGETS
# =========================

.PHONY: all generate clean check-tools

all: generate

## Проверка, что все инструменты установлены
check-tools:
	@command -v $(PROTOC) >/dev/null 2>&1 || (echo "protoc not found"; exit 1)
	@command -v $(PROTOC_GEN_GO) >/dev/null 2>&1 || (echo "protoc-gen-go not found"; exit 1)
	@command -v $(PROTOC_GEN_GO_GRPC) >/dev/null 2>&1 || (echo "protoc-gen-go-grpc not found"; exit 1)

## Генерация Go + gRPC
generate: check-tools
	@echo "Generating Go protobuf code..."
	@rm -rf $(GEN_DIR)
	@mkdir -p $(GEN_DIR)
	@$(PROTOC) \
		-I $(PROTO_DIR) \
		--go_out=$(GO_OUT) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GO_GRPC_OUT) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
	@echo "Done."

## Очистка сгенерированного кода
clean:
	@rm -rf $(GEN_DIR)
	@echo "Generated files removed."