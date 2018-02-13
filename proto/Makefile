PROTOCINC = -I.
PROTOCINC += -Ivendor/github.com/googleapis/googleapis/third_party/protobuf/src/
PROTOCINC += -Ivendor/github.com/googleapis/googleapis/
PROTOCINC += -Ivendor/github.com/googleapis/googleapis/third_party/
PROTOCINC += -Ivendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis

.PHONY: proto
proto:
	@echo compiling Go proto stubs
	@~/Downloads/protoc/bin/protoc ${PROTOCINC} --go_out=plugins=grpc:proto logservice.proto
