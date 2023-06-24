FROM golang:1.19 as buildbase

WORKDIR /go/src/gitlab.com/rarimo/identity/issuer
COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /usr/local/bin/issuer /go/src/gitlab.com/rarimo/identity/issuer

FROM alpine:3.18.2

RUN apk --update add --no-cache musl libstdc++ gcompat libgomp ca-certificates

WORKDIR /

COPY --from=buildbase "/go/pkg/mod/github.com/iden3/wasmer-go@v0.0.1" "/go/pkg/mod/github.com/iden3/wasmer-go@v0.0.1"

COPY --from=buildbase /usr/local/bin/issuer /usr/local/bin/issuer
COPY --from=buildbase /go/src/gitlab.com/rarimo/identity/issuer/circuits /circuits

ENTRYPOINT ["issuer"]
