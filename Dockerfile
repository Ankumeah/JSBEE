FROM golang:1.27.0-alpine3.23@sha256:3747dcba41c8b0db3211fda4db61638b980e17ac5bb3c94460a975a9cfe19395 AS builder

WORKDIR /opt/backend/

COPY backend/go.mod backend/go.sum ./
RUN --mount=type=cache,target=/root/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go mod download

COPY backend/internal/frontend/ backend/internal/frontend/
RUN  --mount=type=cache,target=/root/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go tool templ generate

ARG db_type=sqlite
COPY backend/ .

RUN  --mount=type=cache,target=/root/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go build \
  -ldflags="-s -w" \
  -tags="${db_type},init" \
  -o ./bin/init \
  ./cmd/init/
RUN  --mount=type=cache,target=/root/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go build \
  -ldflags="-s -w" \
  -tags="${db_type}" \
  -o ./bin/api-server \
  ./cmd/api-server/

FROM nginx:alpine@sha256:289decab414250121a93c3f1b8316b9c69906de3a4993757c424cb964169ad42

WORKDIR /opt/nginx/

ENV GIN_MODE=release
COPY --from=builder /opt/backend/bin/init .
COPY --from=builder /opt/backend/bin/api-server .
COPY  reverse_proxy/ .
RUN chmod +x /opt/nginx/run_nginx.sh

RUN mkdir -p /data/public/

COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

ENTRYPOINT [ "/opt/nginx/entrypoint.sh" ]
