#    ChatApp

## Yêu Cầu


### 1. **Golang**

- Phiên bản yêu cầu: **v1.23.6**
- Kiểm tra:
  ``` sh
  go version
  ```

### 2. **Protobuf & gRPC**

- Cài đặt `protoc-gen-go` và `protoc-gen-go-grpc`:
  ```sh
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```
- Kiểm tra:
  ```sh
   protoc --version
  ```

### 3. **MySQL**

- Phiên bản yêu cầu: **8.0**

- Chạy ở **port mặc định: 3306**
- Kiểm tra MySQL đang chạy:

  ```sh
  mysql --version
  ```

## Chạy Dự Án

1. **Tại thư mục root**

   ```sh
   cd .../ChatApp
   ```

2. **Khởi động các service**

   - Chạy **User Service** (port 50052): 
     ```sh
     make start-user
     ```
   - Chạy **Auth Service** (port 50051):
     ```sh
     make start-auth
     ```
   - Chạy **Gateway** (port 8080):
     ```sh
     make start-gateway
     ```

## Kiểm tra gRPC bằng BloomRPC
- Cài đặt: https://github.com/bloomrpc/bloomrpc/releases
- Import file **.proto** và điều chỉnh address port
