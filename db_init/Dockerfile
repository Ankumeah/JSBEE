FROM golang:alpine@sha256:d9b2e14101f27ec8d09674cd01186798d227bb0daec90e032aeb1cd22ac0f029 AS builder

WORKDIR /opt/db_init/

COPY go.mod go.sum ./
RUN go mod download

ARG db_type=sqlite

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
  go build \
  -ldflags="-s -w" \
  -tags=${db_type} \
  -o ./bin/db-init \
  ./

FROM scratch

COPY --from=builder /opt/db_init/bin/db-init .

ENTRYPOINT [ "./db-init" ]
