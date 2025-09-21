

LOCAL_PATH := validate
PROTO_PATH := example
PROTO_FILES := $(shell find example -name "*.proto")

build:
	go build -o protoc-gen-validate main.go
	protoc -I $(PROTO_PATH) -I $(LOCAL_PATH) --go_out=paths=source_relative:example/ $(PROTO_FILES)
	PROTOC_GEN_VALIDATE_DEBUG=true INCLUDE_ONEOFS=true protoc \
		-I $(PROTO_PATH) -I $(LOCAL_PATH) \
		--plugin=protoc-gen-validate=./protoc-gen-validate \
		--validate_out=paths=source_relative,outdir=example:.\
		$(PROTO_FILES)

proto:
	protoc -I . --go_out=paths=source_relative:. validate/validate.proto
