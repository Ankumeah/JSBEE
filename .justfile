set dotenv-load
set dotenv-filename := "env.example"

default: _pre_commit

[linux, macos]
_pre_commit *options:
  #cd lib && go vet {{ options }} ./...
  #cd lib && go fmt ./...

  cd relay && go vet {{ options }} ./...
  cd relay && go fmt ./...

  git add .
