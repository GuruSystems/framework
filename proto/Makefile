PROTOCINC = -I.
PROTOCINC += -Isrc/github.com/googleapis/googleapis/third_party/protobuf/src/
PROTOCINC += -Isrc/github.com/googleapis/googleapis/
PROTOCINC += -Isrc/github.com/googleapis/googleapis/third_party/
PROTOCINC += -Isrc/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis

.PHONY: proto
proto:
	@echo compiling auth/auth.proto
	@~/Downloads/protoc/bin/protoc ${PROTOCINC} --go_out=plugins=grpc:auth auth/auth.proto
	@echo compiling keyvalueserver/keyvalueserver.proto
	@~/Downloads/protoc/bin/protoc ${PROTOCINC} --go_out=plugins=grpc:keyvalueserver keyvalueserver/keyvalueserver.proto
	@echo compiling logging/logservice.proto
	@~/Downloads/protoc/bin/protoc ${PROTOCINC} --go_out=plugins=grpc:logging logging/logservice.proto
	@echo compiling registrar/registrar.proto
	@~/Downloads/protoc/bin/protoc ${PROTOCINC} --go_out=plugins=grpc:registrar registrar/registrar.proto
	@echo compiling softcat-rfc-creator/softcat-rfc-creator.proto
	@~/Downloads/protoc/bin/protoc ${PROTOCINC} --go_out=plugins=grpc:softcat-rfc-creator softcat-rfc-creator/softcat-rfc-creator.proto
