cd proto
protoc -I . -I %GOPATH%\pkg\mod --micro_out=. --go_out=. user.proto
cd ..