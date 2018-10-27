cd /pkg/router/proto
protoc -I ./ router.proto --go_out=plugins=grpc:.