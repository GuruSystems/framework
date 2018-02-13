PROTOCINC = -I.
PROTOCINC += -Igithub.com/googleapis/googleapis/third_party/protobuf/src/
PROTOCINC += -Igithub.com/googleapis/googleapis/
PROTOCINC += -Igithub.com/googleapis/googleapis/third_party/
PROTOCINC += -Igithub.com/grpc-ecosystem/grpc-gateway/third_party/googleapis

.PHONY: proto
proto:
	@echo compiling Go proto stubs
	@~/Downloads/protoc/bin/protoc ${PROTOCINC} --go_out=plugins=grpc:auth auth.auth.proto
	@echo compiling Go proto stubs
	@~/Downloads/protoc/bin/protoc ${PROTOCINC} --go_out=plugins=grpc:keyvalueserver keyvalueserver.proto
	@echo compiling Go proto stubs
	@~/Downloads/protoc/bin/protoc ${PROTOCINC} --go_out=plugins=grpc:logging logservice.proto
	@echo compiling Go proto stubs
	@~/Downloads/protoc/bin/protoc ${PROTOCINC} --go_out=plugins=grpc:registrar registrar.proto
