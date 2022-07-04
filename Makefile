default: check test
test:
	go test -v ./...
check:
	golangci-lint run -v ./...

# Install plugins:
#  go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
#  go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
#  go get -u github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc
proto:
	protoc \
		--go_out=./ \
		--go_opt=paths=source_relative \
		--go-grpc_out=./ \
		--go-grpc_opt=paths=source_relative \
		--openapiv2_out=logtostderr=true:./ \
		--grpc-gateway_out=./ \
		--grpc-gateway_opt=logtostderr=true \
		--grpc-gateway_opt=paths=source_relative \
		--doc_out=./pkg/protobuf/ \
		--doc_opt=markdown,transferBoxApi.md \
		--experimental_allow_proto3_optional \
		./pkg/protobuf/*.proto;
.PHONY: test check
