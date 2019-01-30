protoc pkg/gateway/proto/gateway.proto --go_out=plugins=grpc:.
protoc pkg/gateway/proto/job.proto --go_out=plugins=grpc:.


protoc -I=. -I=%GOPATH%/src -I=%GOPATH%/src/github.com/gogo/protobuf/protobuf --gogo_out=. pkg/gateway/proto/gateway.proto
protoc -I=. -I=%GOPATH%/src -I=%GOPATH%/src/github.com/gogo/protobuf/protobuf --gogo_out=. pkg/gateway/proto/job.proto

