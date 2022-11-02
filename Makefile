create:
	protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        proto/config.proto
	protoc -I . --grpc-gateway_out ./ \
        --grpc-gateway_opt logtostderr=true \
        --grpc-gateway_opt paths=source_relative \
        --grpc-gateway_opt generate_unbound_methods=true \
        proto/config.proto

clean:
	rm proto/*.go

run:
	GOOS=linux go build -a -tags netgo -o rusprofile app/cmd/*.go

