start:
	go run cmd/main.go
gen:
	cd internal && protoc --go_out=proto --go-grpc_out=proto proto/user.proto
	cd internal && protoc --go_out=proto --go-grpc_out=proto proto/auth.proto