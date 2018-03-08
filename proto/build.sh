#!/bin/sh

export PROTO_X="auth"
echo "Compiling ${PROTO_X}.proto"
protoc ${PROTOCINC} --go_out=plugins=grpc:. ${PROTO_X}/${PROTO_X}.proto

export PROTO_X="keyvalueserver"
echo "Compiling ${PROTO_X}.proto"
protoc ${PROTOCINC} --go_out=plugins=grpc:. ${PROTO_X}/${PROTO_X}.proto

export PROTO_X="logservice"
echo "Compiling ${PROTO_X}.proto"
protoc ${PROTOCINC} --go_out=plugins=grpc:. ${PROTO_X}/${PROTO_X}.proto

export PROTO_X="registrar"
echo "Compiling ${PROTO_X}.proto"
protoc ${PROTOCINC} --go_out=plugins=grpc:. ${PROTO_X}/${PROTO_X}.proto
