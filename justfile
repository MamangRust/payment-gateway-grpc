migrate:
    go run cmd/migrate/main.go up

migrate-down:
    go run cmd/migrate/main.go down

generate-proto:
    protoc --proto_path=pkg/proto --go_out=internal/pb --go_opt=paths=source_relative --go-grpc_out=internal/pb --go-grpc_opt=paths=source_relative pkg/proto/*.proto

fmt:
    go fmt ./...

vet:
    go vet ./...

lint:
    golangci-lint run

security-scan: govulncheck gosec

govulncheck:
    govulncheck ./...

gosec:
    gosec ./...

run-client:
    go run cmd/client/main.go

run-server:
    go run cmd/server/main.go

gen:
    mockgen -source=pkg/database/schema/db.go -destination=pkg/database/schema/mocks/mock_db.go -package=mock_db
    mockgen -source=internal/mapper/record/interfaces.go -destination=internal/mapper/record/mocks/mock_recordmapper.go -package=mock_recordmapper
    mockgen -source=internal/service/interfaces.go -destination=internal/service/mocks/mock.go
    mockgen -source=internal/repository/interfaces.go -destination=internal/repository/mocks/mock.go
    mockgen -source=internal/mapper/response/service/interfaces.go -destination=internal/mapper/response/mocks/mock.go
    mockgen -source=internal/mapper/record/interfaces.go -destination=internal/mapper/record/mocks/mock.go
    mockgen -source=internal/mapper/proto/interfaces.go -destination=internal/mapper/proto/mocks/mock.go
    mockgen -source=internal/mapper/response/api/interfaces.go -destination=internal/mapper/response/api/mocks/mock.go
    mockgen -source=internal/pb/auth_grpc.pb.go -destination=internal/pb/mocks/auth_grpc_mock.go
    mockgen -source=internal/pb/card_grpc.pb.go -destination=internal/pb/mocks/card_grpc_mock.go
    mockgen -source=internal/pb/merchant_grpc.pb.go -destination=internal/pb/mocks/merchant_grpc_mock.go
    mockgen -source=internal/pb/saldo_grpc.pb.go -destination=internal/pb/mocks/saldo_grpc_mock.go
    mockgen -source=internal/pb/topup_grpc.pb.go -destination=internal/pb/mocks/topup_grpc_mock.go
    mockgen -source=internal/pb/transaction_grpc.pb.go -destination=internal/pb/mocks/transaction_grpc_mock.go
    mockgen -source=internal/pb/transfer_grpc.pb.go -destination=internal/pb/mocks/transfer_grpc_mock.go
    mockgen -source=internal/pb/user_grpc.pb.go -destination=internal/pb/mocks/user_grpc_mock.go
    mockgen -source=internal/pb/withdraw_grpc.pb.go -destination=internal/pb/mocks/withdraw_grpc_mock.go

test:
    go test -race -covermode=atomic -coverprofile=coverage.txt -v ./tests/...


test-all:
    go test -race -covermode=atomic -coverprofile=coverage.txt -v ./...

hurl:
    hurl --test --variable baseUrl=http://localhost:5000 --glob "tests/hurl/*.hurl"

k6 module type:
    k6 run k6/{{module}}/{{module}}_{{type}}.js

k6-smoke:
    just k6 auth common
    just k6 merchant common
    just k6 saldo common
    just k6 topup common
    just k6 transaction common
    just k6 transfer common
    just k6 withdraw common

coverage:
    go tool cover -html=coverage.txt

sqlc-generate:
    sqlc generate

db-docs:
    dbdocs build doc/db.dbml

db-schema:
    dbml2sql --postgres -o doc/schema.sql doc/db.dbml

db-sqltodbml:
    sql2dbml --postgres doc/schema.sql -o doc/db.dbml

generate-swagger:
    swag init -g cmd/client/main.go

docker-up:
    docker compose up -d --build

docker-down:
    docker compose down

build-client:
    go build -ldflags="-s -w" -o client ./cmd/client

build-server:
    go build -ldflags="-s -w" -o server ./cmd/server
