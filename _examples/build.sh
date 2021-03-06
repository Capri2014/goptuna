#!/bin/sh

export GO111MODULE=on
DIR=$(cd $(dirname $0); pwd)
BIN_DIR=$(cd $(dirname $(dirname $0)); pwd)/bin

mkdir -p ${BIN_DIR}
go build -o ${BIN_DIR}/concurrency ${DIR}/concurrency/main.go
go build -o ${BIN_DIR}/trialnotify ${DIR}/trialnotify/main.go
go build -o ${BIN_DIR}/signalhandling ${DIR}/signalhandling/main.go
go build -o ${BIN_DIR}/simple_rdb ${DIR}/simple_rdb/main.go
go build -o ${BIN_DIR}/simple_tpe ${DIR}/simple_tpe/main.go
