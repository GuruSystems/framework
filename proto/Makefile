PROTOCINC = -I.
PROTOCINC += -Isrc/github.com/googleapis/googleapis/third_party/protobuf/src/
PROTOCINC += -Isrc/github.com/googleapis/googleapis/
PROTOCINC += -Isrc/github.com/googleapis/googleapis/third_party/
PROTOCINC += -Isrc/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis

.PHONY: proto
all:
	auth
	keyvalue
	logging
	registrar
	softcat-rfc
auth:
	@echo compiling auth/auth.proto
	@protoc ${PROTOCINC} --go_out=plugins=grpc:auth auth/auth.proto
keyvalue:
	@echo compiling keyvalueserver/keyvalueserver.proto
	@protoc ${PROTOCINC} --go_out=plugins=grpc:keyvalueserver keyvalueserver/keyvalueserver.proto
logging:
	@echo compiling logging/logservice.proto
	@protoc ${PROTOCINC} --go_out=plugins=grpc:logging logging/logservice.proto
registrar:
	@echo compiling registrar/registrar.proto
	@protoc ${PROTOCINC} --go_out=plugins=grpc:registrar registrar/registrar.proto
softcat-rfc:
	@echo compiling softcat-rfc-creator/softcat-rfc-creator.proto
	@protoc ${PROTOCINC} --go_out=plugins=grpc:softcat-rfc-creator softcat-rfc-creator/softcat-rfc-creator.proto
