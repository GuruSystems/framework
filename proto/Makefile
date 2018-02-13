PROTOCINC = -I.
PROTOCINC += -Isrc/github.com/googleapis/googleapis/third_party/protobuf/src/
PROTOCINC += -Isrc/github.com/googleapis/googleapis/
PROTOCINC += -Isrc/github.com/googleapis/googleapis/third_party/
PROTOCINC += -Isrc/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis

.PHONY: proto
proto:
	@echo compiling Go proto stubs
	@~/Downloads/protoc/bin/protoc ${PROTOCINC} --go_out=plugins=grpc:auth auth/auth.proto
	@echo compiling Go proto stubs
	@~/Downloads/protoc/bin/protoc ${PROTOCINC} --go_out=plugins=grpc:keyvalueserver keyvalueserver/keyvalueserver.proto
	@echo compiling Go proto stubs
	@~/Downloads/protoc/bin/protoc ${PROTOCINC} --go_out=plugins=grpc:logging logging/logservice.proto
	@echo compiling Go proto stubs
	@~/Downloads/protoc/bin/protoc ${PROTOCINC} --go_out=plugins=grpc:registrar registrar/registrar.proto
