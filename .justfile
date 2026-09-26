set dotenv-load
set dotenv-filename := "env.example"

default: _pre_commit

[linux, macos]
test *options:
  cd backend && go test {{ options }} -tags=sqlite,init ./...
  cd smtp_relay && go test {{ options }} ./...

[linux, macos]
up:
  docker compose up --build

[linux, macos]
_pre_commit *options="-tags='sqlite,init'":
  cd backend && go tool templ generate ./internal/frontend/

  cd backend && go vet {{ options }} ./...
  cd backend && go test {{ options }} ./...
  cd backend && go fmt ./...

  cd smtp_relay && go vet {{ options }} ./...
  cd smtp_relay && go test {{ options }} ./...
  cd smtp_relay && go fmt ./...

  git add .
