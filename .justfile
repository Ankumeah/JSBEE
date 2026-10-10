set dotenv-load
set dotenv-filename := "env.example"

default: ci

[linux, macos]
up:
  docker compose up --build

[linux, macos]
ci *options="-tags=sqlite,modernc,init":
  ./.ci.sh {{ options }}

  git add .
