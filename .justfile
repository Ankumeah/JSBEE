set dotenv-load
set dotenv-filename := "env.example"

default: _pre_commit

[linux, macos]
up:
  docker compose up --build

[linux, macos]
_pre_commit *options="-tags=sqlite,modernc,init":
  ./.ci.sh {{ options }}

  git add .
