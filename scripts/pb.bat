@echo off

cd %~dp0
set DIR=%~dp0
set INPUTDIR=%DIR%/../proto/msg
set OUTDIR=%DIR%/../proto

protoc --proto_path=%INPUTDIR%  --gofast_out=%OUTDIR% %INPUTDIR%/*.proto
protoc-go-inject-tag -input=%OUTDIR%/pb/\*.pb.go


echo bingo!