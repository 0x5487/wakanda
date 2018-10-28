cd /pkg/router/proto
protoc -I ./ router.proto --go_out=plugins=grpc:.



protoc pkg/router/proto/router.proto --go_out=plugins=grpc:.