FROM golang:alpine@sha256:d9b2e14101f27ec8d09674cd01186798d227bb0daec90e032aeb1cd22ac0f029 AS builder

WORKDIR /opt/backend/

COPY go.mod go.sum ./
RUN go mod download

ARG db_type=sqlite
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
  go build \
  -ldflags="-s -w" \
  -tags=${db_type} \
  -o ./bin/api-server \
  ./cmd/api-server/

FROM scratch

ENV GIN_MODE=release
COPY --from=builder /opt/backend/bin/api-server .

ENTRYPOINT [ "./api-server" ]
