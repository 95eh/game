#!/usr/bin/env bash

DIR=`dirname $0`
PBDIR=$DIR/../pkg/proto/msg
OUTDIR=$DIR/../pkg/proto

echo PBDIR: $PBDIR
echo OUTDIR: $OUTDIR

echo begin

protoc --proto_path=$PBDIR  --gofast_out=$OUTDIR $PBDIR/*.proto
protoc-go-inject-tag -input=$OUTDIR/pb/\*.pb.go #*要传给命令程序处理，而不是被shell理解为通配符，所以转义

echo finished